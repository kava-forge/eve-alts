package tags

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

func NewTagsTab(deps dependencies, parent fyne.Window) (fyne.CanvasObject, *bindings.DataList[*repository.TagDBData]) {
	logger := logging.With(deps.Logger(), keys.Component, "MainWindow.TagsTab")

	lout := layout.NewGridWrapLayout(fyne.Size{Width: 250, Height: 50})
	ctn := container.New(lout)

	tagsData := bindings.NewDataList[*repository.TagDBData]()
	tagsData.AddListener(binding.NewDataListener(func() {
		defer func() {
			level.Debug(logger).Message("tags listener done")
		}()
		// defer ctn.Refresh()

		level.Debug(logger).Message("tags listener start")

		level.Debug(logger).Message("refreshing tags shown")
		tagsShown := make(map[int64]int, tagsData.Length())
		toDelete := make(map[int64]int, tagsData.Length())

		for i, o := range ctn.Objects {
			if c, ok := o.(*TagCard); ok {
				tagsShown[c.TagID()] = i
				toDelete[c.TagID()] = i
			}
		}

		tagList, err := tagsData.Get()
		if err != nil {
			apperrors.Show(logger, parent, apperrors.Error(
				"Could not find tag list data",
				apperrors.WithCause(err),
			), nil)
			return
		}

		level.Debug(logger).Message("tag list", "shown", tagsShown, "toDelete", toDelete)

		for i, tag := range tagList {
			if tag == nil || tag.Tag.ID == 0 {
				continue
			}

			logger := logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name) //nolint:govet // intentional

			if idx := tagsShown[tag.Tag.ID]; idx != 0 {
				level.Debug(logger).Message("refreshing tag", keys.TagID, tag.Tag.ID, "index", idx)
				delete(toDelete, tag.Tag.ID)
				if err := ctn.Objects[idx].(*TagCard).RefreshDataWith(tag); err != nil {
					apperrors.Show(logger, parent, apperrors.Error(
						"Could not set tags data",
						apperrors.WithCause(err),
					), nil)
				}
				continue
			}

			cc := NewTagCard(deps, parent, tagsData.Child(i), func(c *TagCard) {
				level.Debug(logger).Message("removing tag card")
				tag, err := tagsData.GetValue(i) //nolint:govet // intentional
				if err != nil {
					apperrors.Show(logger, c.Parent(), apperrors.Error(
						"Could not find tag data",
						apperrors.WithCause(err),
					), nil)
					return
				}
				tag.Tag.ID = 0
				if err := tagsData.SetValue(i, tag); err != nil {
					apperrors.Show(logger, c.Parent(), apperrors.Error(
						"Could not delete tag",
						apperrors.WithCause(err),
					), nil)
				}
				ctn.Remove(c)
			}, func(tagData bindings.DataProxy[*repository.TagDBData], onClose func()) {
				w := NewTagEditor(deps, fyne.CurrentApp(), "Edit Tag", tagsData, tagData, onClose)
				w.Show()
			})
			level.Debug(logger).Message("adding new tag", keys.TagID, tag.Tag.ID, "cc", fmt.Sprintf("%#v", cc))
			ctn.Add(cc)
		}

		for tagID, idx := range toDelete {
			level.Debug(logger).Message("removing widget", keys.TagID, tagID, "index", idx)
			ctn.Remove(ctn.Objects[idx])
		}
	}))

	ctn.Add(NewAddTagButton(deps, tagsData))

	return container.NewVScroll(ctn), tagsData
}
