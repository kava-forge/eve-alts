package characters

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func NewCharactersTab(deps dependencies, parent fyne.Window, tags *bindings.DataList[*repository.TagDBData]) (fyne.CanvasObject, *bindings.DataList[*repository.CharacterDBData]) {
	logger := logging.With(deps.Logger(), keys.Component, "MainWindow.CharactersTab")

	tagFilters := NewTagFilter(deps, parent, tags)

	charLout := layout.NewGridWrapLayout(fyne.Size{Width: 500, Height: 128})
	charContainer := container.New(charLout)

	chars := bindings.NewDataList[*repository.CharacterDBData]()
	chars.AddListener(binding.NewDataListener(func() {
		defer charContainer.Refresh()

		level.Debug(logger).Message("refreshing characters shown")
		charsShown := make(map[int64]int, chars.Length())
		toDelete := make(map[int64]int, chars.Length())

		for i, o := range charContainer.Objects {
			if c, ok := o.(*CharacterCard); ok {
				charsShown[c.CharacterID()] = i
				toDelete[c.CharacterID()] = i
			}
		}

		charList, err := chars.Get()
		if err != nil {
			apperrors.Show(logger, parent, apperrors.Error(
				"Could not find character list data",
				apperrors.WithCause(err),
			), nil)
			return
		}

		level.Debug(logger).Message("char list", "shown", charsShown, "toDelete", toDelete)

		for i, char := range charList {
			if char.Character.ID == 0 {
				continue
			}

			logger := logging.With(logger, keys.CharacterID, char.Character.ID, keys.CharacterName, char.Character.Name) //nolint:govet // intentional

			if idx := charsShown[char.Character.ID]; idx != 0 {
				level.Debug(logger).Message("refreshing character", "index", idx)
				delete(toDelete, char.Character.ID)
				if err := charContainer.Objects[idx].(*CharacterCard).RefreshDataWith(*char); err != nil {
					apperrors.Show(logger, parent, apperrors.Error(
						"Could not set character data",
						apperrors.WithCause(err),
					), nil)
				}
				continue
			}

			cc := NewCharacterCard(deps, parent, chars.Child(i), tags, func(c *CharacterCard) {
				level.Debug(logger).Message("removing character card")
				char, err := chars.GetValue(i)
				if err != nil {
					apperrors.Show(logger, c.Parent(), apperrors.Error(
						"Could not find character data",
						apperrors.WithCause(err),
					), nil)
					return
				}
				char.Character.ID = 0
				if err := chars.SetValue(i, char); err != nil {
					apperrors.Show(logger, c.Parent(), apperrors.Error(
						"Could not delete character",
						apperrors.WithCause(err),
					), nil)
					return
				}
				charContainer.Remove(c)
			})
			level.Debug(logger).Message("adding new character", "cc", fmt.Sprintf("%#v", cc))
			charContainer.Add(cc)

			tagFilters.AttachAllToCharacter(cc)
		}

		for charID, idx := range toDelete {
			level.Debug(logger).Message("removing widget", keys.CharacterID, charID, "index", idx)
			charContainer.Remove(charContainer.Objects[idx])
		}
	}))

	buttonLout := layout.NewGridWrapLayout(fyne.Size{Width: 500, Height: 50})
	buttonContainer := container.New(buttonLout)

	buttonContainer.Add(NewAddCharacterButton(deps, parent, chars))
	buttonContainer.Add(NewRefreshAllButton(deps, parent, chars))

	vbox := container.New(layout.NewVBoxLayout())
	vbox.Add(buttonContainer)
	vbox.Add(container.New(layout.NewPaddedLayout(), tagFilters.TagSet))
	vbox.Add(charContainer)

	return container.NewVScroll(vbox), chars
}
