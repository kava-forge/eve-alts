package characters

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/panics"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func NewAddCharacterButton(deps dependencies, parent fyne.Window, chars *bindings.DataList[*repository.CharacterDBData]) *widget.Button {
	button := widget.NewButtonWithIcon("Add Character", theme.ContentAddIcon(), func() {
		go func() {
			logger := logging.With(deps.Logger(), keys.Component, "AddCharacterButton")

			defer panics.Handler(logger)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			tok, err := deps.ESIClient().Authenticate(ctx)
			if err != nil {
				apperrors.Show(logger, parent, apperrors.Error(
					"Could not authenticate with ESI",
					apperrors.WithCause(err),
				), nil)
				return
			}

			cdata, err := deps.ESIClient().ValidateToken(ctx, tok)
			if err != nil {
				apperrors.Show(logger, parent, apperrors.Error(
					"ESI Token invalid. Please re-add the character",
					apperrors.WithCause(err),
				), nil)
				return
			}

			dbChar, err := RefreshCharacterData(ctx, deps, tok, cdata.RealID)
			if err != nil {
				apperrors.Show(logger, parent, apperrors.Error(
					"Error fetching character data",
					apperrors.WithCause(err),
				), nil)
				return
			}

			if err := chars.Append(&dbChar); err != nil {
				apperrors.Show(logger, parent, apperrors.Error(
					"Error displaying character",
					apperrors.WithCause(err),
					apperrors.WithInternalData(keys.CharacterID, dbChar.Character.ID),
				), nil)
				return
			}
		}()
	})

	return button
}
