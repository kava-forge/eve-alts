package characters

import (
	"context"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
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
		pb := dialog.NewCustomWithoutButtons("Refreshing All Characters", widget.NewProgressBarInfinite(), parent)
		pb.Show()
		defer pb.Hide()

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

		var me error
		errs := make(chan error, len(charList))
		done := make(chan struct{})
		go func() {
			for err := range errs {
				me = multierror.Append(me, err)
			}
			close(done)
		}()

		wg := &sync.WaitGroup{}

		for i, char := range charList {
			if char == nil || char.Character.ID == 0 {
				continue
			}

			char := char
			wg.Add(1)
			go func() {
				defer wg.Done()

				dbTok, err := deps.AppRepo().GetTokenForCharacter(ctx, char.Character.ID, nil)
				if err != nil {
					errs <- errors.Wrap(err, "unable to find db token", keys.CharacterID, char.Character.ID)
					return
				}

				tok := esi.TokenFromRepository(dbTok)

				dbChar, err := RefreshCharacterData(ctx, deps, tok, char.Character.ID)
				if err != nil {
					errs <- errors.Wrap(err, "could not RefreshCharacterData", keys.CharacterID, char.Character.ID)
					return
				}

				if err := chars.SetValue(i, &dbChar); err != nil {
					errs <- errors.Wrap(err, "could not append new character", keys.CharacterID, char.Character.ID)
				}
			}()
		}

		wg.Wait()
		close(errs)
		<-done

		if me != nil {
			apperrors.Show(logger, parent, apperrors.Error(
				"Could not refresh all characters",
				apperrors.WithCause(me),
			), nil)
		}
	}

	return button
}
