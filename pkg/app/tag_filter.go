package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type TagFilter struct {
	deps               dependencies
	parent             fyne.Window
	tags               *DataList[*repository.TagDBData]
	TagSet             *MiniTagSet[string, *TagMiniTag]
	tagState           *DataMap[bool]
	attachedCharacters []*CharacterCard
}

func NewTagFilter(deps dependencies, parent fyne.Window, tags *DataList[*repository.TagDBData]) *TagFilter {
	tf := &TagFilter{
		deps:   deps,
		parent: parent,
		tags:   tags,

		TagSet:   NewMiniTagSet[string, *TagMiniTag](),
		tagState: NewDataMap[bool](),
	}

	tf.update()

	tags.AddListener(binding.NewDataListener(tf.update))

	return tf
}

func (tf *TagFilter) update() {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.tagsListener")

	level.Debug(logger).Message("handling filter update")

	tagsList, err := tf.tags.Get()
	if err != nil {
		ShowError(logger, tf.parent, AppError(
			"Could not load tag list data",
			WithCause(err),
		), nil)
		return
	}

	for i, t := range tagsList {
		if tf.tagState.HasKey(t.StrID()) {
			level.Debug(logger).Message("skipping because key exists", "tag_id", t.StrID())
			continue
		}

		tf.Add(tf.tags.Child(i))

		for _, cc := range tf.attachedCharacters {
			tf.AttachToCharacter(cc, t.StrID(), tf.tagState.Child(t.StrID()))
		}
	}
}

// func (tf *TagFilter) resetIfAllOff() {
// 	defer tf.TagSet.Refresh()

// 	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.resetIfAllOff")

// 	data, err := tf.tagState.Get()
// 	if err != nil {
// 		ShowError(logger, tf.parent, AppError(
// 			"Could not load tag data",
// 			WithCause(err),
// 		), nil)
// 		return
// 	}

// 	anyOn := false

// 	for _, v := range data {
// 		if v {
// 			anyOn = true
// 		}
// 	}

// 	if !anyOn {
// 		for k := range data {
// 			tf.tagState.SetValue(k, true)
// 		}
// 	}
// }

func (tf *TagFilter) Add(tagData DataProxy[*repository.TagDBData]) {
	defer tf.TagSet.Refresh()

	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.Add")

	tag, err := tagData.Get()
	if err != nil {
		ShowError(logger, tf.parent, AppError(
			"Could not load tag data",
			WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	if err := tf.tagState.SetValue(tag.StrID(), false); err != nil {
		ShowError(logger, tf.parent, AppError(
			"Could not set filter state",
			WithCause(err),
		), nil)
		return
	}

	level.Debug(logger).Message("adding mini tag")

	stateChild := tf.tagState.Child(tag.StrID())
	tmt := NewTagMiniTag(tf.deps, tf.parent, tagData, stateChild)
	tf.TagSet.Add(tmt)

	// tf.resetIfAllOff()
	// stateChild.AddListener(binding.NewDataListener(tf.resetIfAllOff))
}

func (tf *TagFilter) AttachAllToCharacter(cc *CharacterCard) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.AttachAllToCharacter")

	tagIDs := tf.tagState.Keys()

	level.Debug(logger).Message("attaching all to character", "tag_ids", tagIDs)

	for _, tagID := range tagIDs {
		tagSelected := tf.tagState.Child(tagID)
		tf.AttachToCharacter(cc, tagID, tagSelected)
	}

	tf.attachedCharacters = append(tf.attachedCharacters, cc)
}

func (tf *TagFilter) AttachToCharacter(cc *CharacterCard, tagID string, tagSelected DataProxy[bool]) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.AttachToCharacter")

	level.Info(logger).Message("attaching to character", "tag_id", tagID)

	tagSelected.AddListener(binding.NewDataListener(func() {
		level.Debug(logger).Message("refreshing character for filter change")
		selected, err := tagSelected.Get()
		if err != nil {
			ShowError(logger, tf.parent, AppError(
				"Could not get filter state",
				WithCause(err),
			), nil)
			return
		}

		cc.selectedTags[tagID] = selected
		cc.redraw()
	}))
}
