package roles

import (
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/hashicorp/go-multierror"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/minitag"
	"github.com/kava-forge/eve-alts/pkg/app/tags"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type TagSelector struct {
	deps     dependencies
	parent   fyne.Window
	tags     *bindings.DataList[*repository.TagDBData]
	TagSet   *minitag.MiniTagSet[string, *tags.TagMiniTag]
	tagState *bindings.DataMap[bool]
}

func NewTagSelector(deps dependencies, parent fyne.Window, tagsData *bindings.DataList[*repository.TagDBData]) *TagSelector {
	tf := &TagSelector{
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

func (tf *TagSelector) update() {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagSelector.tagsListener")

	level.Debug(logger).Message("handling selector update")

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
	}
	tf.TagSet.Refresh()
}

func (tf *TagSelector) Selected() []int64 {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagSelector.Selected")

	tagState, err := tf.tagState.Get()
	if err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not load tag selection data",
			apperrors.WithCause(err),
		), nil)
		return nil
	}

	sel := make([]int64, 0, len(tagState))
	for k, v := range tagState {
		if !v {
			continue
		}

		ki, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			level.Error(logger).Err("could not parse selection key", err)
			continue
		}

		sel = append(sel, ki)
	}

	slices.Sort(sel)
	return sel
}

func (tf *TagSelector) Select(tagID int64) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagSelector.Selected")

	tagIDStr := strconv.FormatInt(tagID, 10)
	if err := tf.tagState.SetValue(tagIDStr, true); err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not set tag selection",
			apperrors.WithCause(err),
		), nil)
	}
}

func (tf *TagSelector) Clear() {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagSelector.Clear")

	var errs error

	for _, k := range tf.tagState.Keys() {
		if err := tf.tagState.SetValue(k, false); err != nil {
			errs = multierror.Append(errs, errors.WithDetails(err, keys.TagID, k))
		}
	}

	if errs != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not clear tag selections",
			apperrors.WithCause(errs),
		), nil)
	}
}

func (tf *TagSelector) add(tagData bindings.DataProxy[*repository.TagDBData]) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "TagSelector.Add")

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
