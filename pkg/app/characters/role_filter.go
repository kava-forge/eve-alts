package characters

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/minitag"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type RoleFilter struct {
	deps               dependencies
	parent             fyne.Window
	roles              *bindings.DataList[*repository.RoleDBData]
	RoleSet            *minitag.MiniTagSet[string, *RoleFilterTag]
	roleState          *bindings.DataMap[bool]
	attachedCharacters []*CharacterCard
}

func NewRoleFilter(deps dependencies, parent fyne.Window, rolesData *bindings.DataList[*repository.RoleDBData]) *RoleFilter {
	tf := &RoleFilter{
		deps:   deps,
		parent: parent,
		roles:  rolesData,

		RoleSet:   minitag.NewMiniTagSet[string, *RoleFilterTag](),
		roleState: bindings.NewDataMap[bool](),
	}

	tf.update()

	rolesData.AddListener(binding.NewDataListener(tf.update))

	return tf
}

func (tf *RoleFilter) update() {
	logger := logging.With(tf.deps.Logger(), keys.Component, "RoleFilter.update")

	level.Debug(logger).Message("handling filter update")

	rolesList, err := tf.roles.Get()
	if err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not load role list data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	for i, r := range rolesList {
		if tf.roleState.HasKey(r.StrID()) {
			level.Debug(logger).Message("skipping because key exists", "role_id", r.StrID())
			continue
		}

		tf.add(tf.roles.Child(i))

		for _, cc := range tf.attachedCharacters {
			tf.AttachToCharacter(cc, r.StrID(), tf.roleState.Child(r.StrID()))
		}
	}
	tf.RoleSet.Refresh()
}

func (tf *RoleFilter) add(roleData bindings.DataProxy[*repository.RoleDBData]) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "RoleFilter.Add")

	role, err := roleData.Get()
	if err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not load role data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.RoleID, role.Role.ID, keys.RoleName, role.Role.Name)

	if err := tf.roleState.SetValue(role.StrID(), false); err != nil {
		apperrors.Show(logger, tf.parent, apperrors.Error(
			"Could not set filter state",
			apperrors.WithCause(err),
		), nil)
		return
	}

	level.Debug(logger).Message("adding mini tag")

	stateChild := tf.roleState.Child(role.StrID())
	tmt := NewRoleFilterTag(tf.deps, tf.parent, roleData, stateChild)
	tf.RoleSet.Add(tmt)
}

func (tf *RoleFilter) AttachAllToCharacter(cc *CharacterCard) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "RoleFilter.AttachAllToCharacter")

	roleIDs := tf.roleState.Keys()

	level.Debug(logger).Message("attaching all to character", "role_ids", roleIDs)

	for _, roleID := range roleIDs {
		roleSelected := tf.roleState.Child(roleID)
		tf.AttachToCharacter(cc, roleID, roleSelected)
	}

	tf.attachedCharacters = append(tf.attachedCharacters, cc)
}

func (tf *RoleFilter) AttachToCharacter(cc *CharacterCard, roleID string, roleSelected bindings.DataProxy[bool]) {
	logger := logging.With(tf.deps.Logger(), keys.Component, "RoleFilter.AttachToCharacter")

	level.Info(logger).Message("attaching to character", "role_id", roleID)

	roleSelected.AddListener(bindings.NewListener(logger, func() {
		level.Debug(logger).Message("refreshing character for role filter change")
		selected, err := roleSelected.Get()
		if err != nil {
			apperrors.Show(logger, tf.parent, apperrors.Error(
				"Could not get role filter state",
				apperrors.WithCause(err),
			), nil)
			return
		}

		cc.UpdateRoleSelection(roleID, selected)
	}))
}
