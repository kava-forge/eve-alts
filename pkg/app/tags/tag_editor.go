package tags

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"image/color"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/hashicorp/go-multierror"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/colors"
	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func populateTagData(deps dependencies, nameInp *widget.Entry, colorSwatch *colors.TappableColorSwatch, colorInp *dialog.ColorPickerDialog, textArea *widget.Entry, tagData bindings.DataProxy[*repository.TagDBData]) error {
	tag, err := tagData.Get()
	if err != nil {
		return errors.Wrap(err, "unable to load tag data")
	}

	nameInp.Text = tag.Tag.Name
	nameInp.Refresh()
	colorSwatch.SetColor(tag.Color())
	colorSwatch.Refresh()
	colorInp.SetColor(tag.Color())

	lines := make([]string, 0, len(tag.Skills))

	skillIDs := make([]int64, 0, len(tag.Skills))
	for _, sk := range tag.Skills {
		skillIDs = append(skillIDs, sk.SkillID)
	}
	rows, err := deps.StaticRepo().BatchGetSkillNames(context.Background(), skillIDs, nil)
	if err != nil {
		return errors.Wrap(err, "could not fetch skill names")
	}
	nameMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		nameMap[row.SkillID] = row.SkillName
	}

	for _, sk := range tag.Skills {
		lines = append(lines, fmt.Sprintf("%s %d", nameMap[sk.SkillID], sk.SkillLevel))
	}
	sort.Strings(lines)

	textArea.Text = strings.Join(lines, "\n")
	textArea.Refresh()

	return nil
}

func NewTagEditor(deps dependencies, a fyne.App, title string, tags *bindings.DataList[*repository.TagDBData], tagData bindings.DataProxy[*repository.TagDBData], onClose func()) fyne.Window {
	logger := logging.With(deps.Logger(), keys.Component, "TagEditor")

	w := a.NewWindow(fmt.Sprintf("EVE Alts - %s", title))
	w.CenterOnScreen()
	w.Resize(fyne.Size{Width: 800, Height: 600})
	w.SetOnClosed(func() {
		if onClose != nil {
			onClose()
		}
	})

	nameInp := widget.NewEntry()

	colorSwatch := colors.NewTappableColorSwatch(deps.Logger(), color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}) //nolint:gosec // not security related
	colorInp := dialog.NewColorPicker("Tag Color", "Pick a color for the tag", func(c color.Color) {
		colorSwatch.SetColor(c)
	}, w)
	colorInp.Advanced = true
	colorSwatch.OnTapped = func(pe *fyne.PointEvent) { colorInp.Show() }

	textArea := widget.NewMultiLineEntry()
	textArea.Wrapping = fyne.TextWrapOff
	textArea.SetMinRowsVisible(10)

	if tagData != nil {
		if err := populateTagData(deps, nameInp, colorSwatch, colorInp, textArea, tagData); err != nil {
			apperrors.Show(logger, w, apperrors.Error(
				"Could not load tag data",
				apperrors.WithCause(err),
			), nil)
		}
	}

	form := widget.NewForm(
		widget.NewFormItem("Tag Name", nameInp),
		widget.NewFormItem("Tag Color", colorSwatch),
		widget.NewFormItem("Skill List", textArea),
	)
	form.OnCancel = w.Close
	form.OnSubmit = func() {
		ctx := context.Background()

		skills, err := ParseSkills(ctx, deps.StaticRepo(), textArea.Text)
		if err != nil {
			apperrors.Show(logger, w, apperrors.Error(
				"Could not parse skill list",
				apperrors.WithCause(err),
			), nil)
			return
		}
		level.Debug(logger).Message("parsed skills", "skills", skills)

		if tagData == nil {
			var dbTag repository.Tag
			var dbSkills []repository.TagSkill
			if err := database.TransactWithRetries(ctx, deps.Telemetry(), logger, deps.DB(), &sql.TxOptions{}, func(ctx context.Context, tx database.Tx) error {
				if dbSkills == nil {
					dbSkills = make([]repository.TagSkill, 0, len(skills))
				} else {
					dbSkills = dbSkills[:0]
				}

				dbTag, err = deps.AppRepo().InsertTag(ctx, nameInp.Text, colorSwatch.Color(), tx)
				if err != nil {
					return errors.Wrap(err, "could not InsertTag")
				}

				for _, sd := range skills {
					sk, err := deps.AppRepo().UpsertTagSkill(ctx, dbTag.ID, sd.SkillID, sd.SkillLevel, tx)
					if err != nil {
						return errors.Wrap(err, "could not UpsertTagSkill", keys.SkillID, sd.SkillID, keys.SkillLevel, sd.SkillLevel)
					}
					dbSkills = append(dbSkills, sk)
				}

				return nil
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not create tag",
					apperrors.WithCause(err),
				), nil)
				return
			}

			if err := tags.Append(&repository.TagDBData{
				Tag:    dbTag,
				Skills: dbSkills,
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not append tag",
					apperrors.WithCause(err),
				), nil)
			}
			w.Close()
		} else { // edit existing
			tagP, err := tagData.Get()
			if err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not find tag data",
					apperrors.WithCause(err),
				), nil)
				return
			}

			logger := logging.With(logger, keys.TagID, tagP.Tag.ID) //nolint:govet // intentional

			var dbTag repository.Tag
			var dbSkills []repository.TagSkill
			if err := database.TransactWithRetries(ctx, deps.Telemetry(), logger, deps.DB(), &sql.TxOptions{}, func(ctx context.Context, tx database.Tx) error {
				tag := *tagP

				skillIDMap := make(map[int64]bool, len(skills))
				for _, skill := range skills {
					skillIDMap[skill.SkillID] = true
				}
				toDelete := make([]int64, 0, 10)
				for _, skill := range tag.Skills {
					if !skillIDMap[skill.SkillID] {
						toDelete = append(toDelete, skill.SkillID)
					}
				}

				if dbSkills == nil {
					dbSkills = make([]repository.TagSkill, 0, len(skills))
				} else {
					dbSkills = dbSkills[:0]
				}

				c := colorSwatch.Color()
				err = deps.AppRepo().UpdateTag(ctx, tag.Tag.ID, nameInp.Text, c, tx)
				if err != nil {
					return errors.Wrap(err, "could not UpdateTag")
				}
				tag.Tag.Name = nameInp.Text
				cr, cg, cb, _ := c.RGBA()
				tag.Tag.ColorR = int64(cr)
				tag.Tag.ColorG = int64(cg)
				tag.Tag.ColorB = int64(cb)

				for _, sd := range skills {
					sk, err := deps.AppRepo().UpsertTagSkill(ctx, tag.Tag.ID, sd.SkillID, sd.SkillLevel, tx)
					if err != nil {
						return errors.Wrap(err, "could not UpsertTagSkill", keys.SkillID, sd.SkillID, keys.SkillLevel, sd.SkillLevel)
					}
					dbSkills = append(dbSkills, sk)
				}

				if len(toDelete) > 0 {
					if err := deps.AppRepo().DeleteTagSkills(ctx, tag.Tag.ID, toDelete, tx); err != nil {
						return errors.Wrap(err, "could not DeleteTagSkills")
					}
				}

				dbTag = tag.Tag

				return nil
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not save tag",
					apperrors.WithCause(err),
				), nil)
				return
			}

			if err := tagData.Set(&repository.TagDBData{
				Tag:    dbTag,
				Skills: dbSkills,
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not set tag",
					apperrors.WithCause(err),
				), nil)
			}
			w.Close()
		}
	}

	w.SetContent(form)

	return w
}

type SkillData struct {
	SkillID    int64
	SkillLevel int64
}

var (
	skillRegexp         = regexp.MustCompile(`^(.*\D)(\s\d)?$`)
	ErrMalformedLine    = errors.New("malformed line")
	ErrInvalidSkillName = errors.New("invalid skill name")
)

func ParseSkills(ctx context.Context, static repository.StaticData, text string) (skills []SkillData, err error) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	var errs error

	var line string
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		m := skillRegexp.FindStringSubmatch(line)
		if m == nil {
			errs = multierror.Append(errs, errors.Wrap(ErrMalformedLine, "could not parse line", "line", line))
			continue
		}

		skillName := strings.TrimSpace(m[1])
		skillLevelStr := strings.TrimSpace(m[2])

		var skillLevel int64
		if skillLevelStr == "" {
			skillLevel = 1
		} else {
			skillLevel, err = strconv.ParseInt(skillLevelStr, 10, 64)
			if err != nil {
				errs = multierror.Append(errs, errors.Wrap(err, "could not convert skill level to int64", "level", skillLevelStr))
				continue
			}
		}

		skillID, err := static.GetSkillIDByName(ctx, skillName, nil)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrap(ErrInvalidSkillName, "could not find skill", "skill_name", skillName))
			continue
		}

		skills = append(skills, SkillData{
			SkillID:    skillID,
			SkillLevel: skillLevel,
		})
	}

	maxLevels := make(map[int64]SkillData, len(skills))
	for _, sk := range skills {
		if sk.SkillLevel >= maxLevels[sk.SkillID].SkillLevel {
			maxLevels[sk.SkillID] = sk
		}
	}

	skills = skills[:0]
	for _, sk := range maxLevels {
		skills = append(skills, sk)
	}

	return skills, errs
}
