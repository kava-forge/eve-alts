package characters

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hashicorp/go-multierror"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/esi"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func NewRefreshAllButton(deps dependencies, parent fyne.Window, chars *bindings.DataList[*repository.CharacterDBData]) *widget.Button {
	button := widget.NewButtonWithIcon("Refresh All Characters", theme.ViewRefreshIcon(), nil)
	button.OnTapped = func() {
		ctx := context.Background()

		logger := logging.With(deps.Logger(), keys.Component, "RefreshAllButton")

		button.Disable()
		defer button.Enable()

		charList, err := chars.Get()
		if err != nil {
			apperrors.Show(logger, parent, apperrors.Error(
				"Could not find character list data",
				apperrors.WithCause(err),
			), nil)
			return
		}

		var errs error
		for i, char := range charList {
			if char == nil || char.Character.ID == 0 {
				continue
			}

			dbTok, err := deps.AppRepo().GetTokenForCharacter(ctx, char.Character.ID, nil)
			if err != nil {
				errs = multierror.Append(errs, errors.Wrap(err, "unable to find db token", keys.CharacterID, char.Character.ID))
				continue
			}

			tok := esi.TokenFromRepository(dbTok)

			dbChar, err := RefreshCharacterData(ctx, deps, tok, char.Character.ID)
			if err != nil {
				errs = multierror.Append(errs, errors.Wrap(err, "could not RefreshCharacterData", keys.CharacterID, char.Character.ID))
				continue
			}

			if err := chars.SetValue(i, &dbChar); err != nil {
				errs = multierror.Append(errs, errors.Wrap(err, "could not append new character", keys.CharacterID, char.Character.ID))
				continue
			}
		}

		if errs != nil {
			apperrors.Show(logger, parent, apperrors.Error(
				"Could not refresh all characters",
				apperrors.WithCause(err),
			), nil)
		}
	}

	return button
}
