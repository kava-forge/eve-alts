package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func NewTagsTab(deps dependencies, parent fyne.Window) (fyne.CanvasObject, *DataList[*repository.TagDBData]) {
	logger := logging.With(deps.Logger(), keys.Component, "MainWindow.TagsTab")

	lout := layout.NewGridWrapLayout(fyne.Size{Width: 250, Height: 50})
	ctn := container.New(lout)

	tags := NewDataList[*repository.TagDBData]()
	tags.AddListener(binding.NewDataListener(func() {
		defer func() {
			level.Debug(logger).Message("tags listener done")
		}()
		// defer ctn.Refresh()

		level.Debug(logger).Message("tags listener start")

		level.Debug(logger).Message("refreshing tags shown")
		tagsShown := make(map[int64]int, tags.Length())
		toDelete := make(map[int64]int, tags.Length())

		for i, o := range ctn.Objects {
			if c, ok := o.(*TagCard); ok {
				tagsShown[c.TagID()] = i
				toDelete[c.TagID()] = i
			}
		}

		tagList, err := tags.Get()
		if err != nil {
			ShowError(logger, parent, AppError(
				"Could not find tag list data",
				WithCause(err),
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
				if err := ctn.Objects[idx].(*TagCard).RefreshDataWith(*tag); err != nil {
					ShowError(logger, parent, AppError(
						"Could not set tags data",
						WithCause(err),
					), nil)
				}
				continue
			}

			cc := NewTagCard(deps, parent, tags.Child(i), func(c *TagCard) {
				level.Debug(logger).Message("removing tag card")
				tag, err := tags.GetValue(i) //nolint:govet // intentional
				if err != nil {
					ShowError(logger, c.parent, AppError(
						"Could not find tag data",
						WithCause(err),
					), nil)
					return
				}
				tag.Tag.ID = 0
				if err := tags.SetValue(i, tag); err != nil {
					ShowError(logger, c.parent, AppError(
						"Could not delete tag",
						WithCause(err),
					), nil)
				}
				ctn.Remove(c)
			}, func(tagData DataProxy[*repository.TagDBData], onClose func()) {
				w := NewTagEditor(deps, fyne.CurrentApp(), "Edit Tag", tags, tagData, onClose)
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

	ctn.Add(NewAddTagButton(deps, tags))

	return container.NewVScroll(ctn), tags
}
