package characters

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/minitag"
	"github.com/kava-forge/eve-alts/pkg/app/tags"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type TagFilter struct {
	deps               dependencies
	parent             fyne.Window
	tags               *bindings.DataList[*repository.TagDBData]
	TagSet             *minitag.MiniTagSet[string, *tags.TagMiniTag]
	tagState           *bindings.DataMap[bool]
	attachedCharacters []*CharacterCard
}

func NewTagFilter(deps dependencies, parent fyne.Window, tagsData *bindings.DataList[*repository.TagDBData]) *TagFilter {
	tf := &TagFilter{
		deps:   deps,
		parent: parent,
		tags:   tagsData,

		TagSet:   minitag.NewMiniTagSet[string, *tags.TagMiniTag](),
		tagState: bindings.NewDataMap[bool](),
	}

	tf.update()

	tagsData.AddListener(binding.NewDataListener(tf.update))

	return tf
}

func (tf *TagFilter) update() {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.tagsListener")

	level.Debug(logger).Message("handling filter update")

	tagsList, err := tf.tags.Get()
	if err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not load tag list data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	for i, t := range tagsList {
		if tf.tagState.HasKey(t.StrID()) {
			level.Debug(logger).Message("skipping because key exists", "tag_id", t.StrID())
			continue
		}

		tf.add(tf.tags.Child(i))

		for _, cc := range tf.attachedCharacters {
			tf.AttachToCharacter(cc, t.StrID(), tf.tagState.Child(t.StrID()))
		}
	}
	tf.TagSet.Refresh()
}

func (tf *TagFilter) add(tagData bindings.DataProxy[*repository.TagDBData]) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.Add")

	tag, err := tagData.Get()
	if err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not load tag data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	if err := tf.tagState.SetValue(tag.StrID(), false); err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not set filter state",
			apperrors.WithCause(err),
		), nil)
		return
	}

	level.Debug(logger).Message("adding mini tag")

	stateChild := tf.tagState.Child(tag.StrID())
	tmt := tags.NewTagMiniTag(tf.deps, tf.parent, tagData, stateChild)
	tf.TagSet.Add(tmt)
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

func (tf *TagFilter) AttachToCharacter(cc *CharacterCard, tagID string, tagSelected bindings.DataProxy[bool]) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagFilter.AttachToCharacter")

	level.Info(logger).Message("attaching to character", "tag_id", tagID)

	tagSelected.AddListener(bindings.NewListener(func() {
		level.Debug(logger).Message("refreshing character for filter change")
		selected, err := tagSelected.Get()
		if err != nil {
			apperrors.Show(logger, tf.parent, apperrors.Error(
				"Could not get filter state",
				apperrors.WithCause(err),
			), nil)
			return
		}

		cc.UpdateTagSelection(tagID, selected)
	}))
}
