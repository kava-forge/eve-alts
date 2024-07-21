package app

import (
	"context"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/characters"
	"github.com/kava-forge/eve-alts/pkg/app/roles"
	"github.com/kava-forge/eve-alts/pkg/app/tags"
	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/keys"
)

func NewMainWindow(deps dependencies, a fyne.App) fyne.Window {
	logger := logging.With(deps.Logger(), keys.Component, "MainWindow")

	w := a.NewWindow("EVE Alts")
	w.CenterOnScreen()
	w.Resize(fyne.Size{Width: 1024, Height: 768})

	tagsTab, tags := tags.NewTagsTab(deps, w)
	rolesTab, roles := roles.NewRolesTab(deps, w, tags)
	charsTab, chars := characters.NewCharactersTab(deps, w, tags, roles)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Characters", theme.AccountIcon(), charsTab),
		container.NewTabItemWithIcon("Tags", theme.ListIcon(), tagsTab),
		container.NewTabItemWithIcon("Roles", theme.GridIcon(), rolesTab),
	)

	w.SetContent(tabs)

	a.Lifecycle().SetOnStarted(func() {
		ctx := context.Background()

		wg := &sync.WaitGroup{}

		wg.Add(3)

		go func() {
			defer wg.Done()
			knownChars, err := deps.AppRepo().GetAllCharacters(ctx, nil)
			if err != nil && !errors.Is(err, database.ErrNoRows) {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not load character data",
					apperrors.WithCause(err),
				), nil)
			} else {
				level.Debug(logger).Message("initial characters loaded")
				if err := chars.Set(knownChars); err != nil {
					apperrors.Show(logger, w, apperrors.Error(
						"Could not set character list data",
						apperrors.WithCause(err),
					), nil)
				}
			}
		}()

		go func() {
			defer wg.Done()
			knownTags, err := deps.AppRepo().GetAllTags(ctx, nil)
			if err != nil && !errors.Is(err, database.ErrNoRows) {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not load tags data",
					apperrors.WithCause(err),
				), nil)
			} else {
				level.Debug(logger).Message("initial tags loaded")
				if err := tags.Set(knownTags); err != nil {
					apperrors.Show(logger, w, apperrors.Error(
						"Could not set tags list data",
						apperrors.WithCause(err),
					), nil)
				}
			}
		}()

		go func() {
			defer wg.Done()
			knownRoles, err := deps.AppRepo().GetAllRoles(ctx, nil)
			if err != nil && !errors.Is(err, database.ErrNoRows) {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not load roles data",
					apperrors.WithCause(err),
				), nil)
			} else {
				level.Debug(logger).Message("initial roles loaded")
				if err := roles.Set(knownRoles); err != nil {
					apperrors.Show(logger, w, apperrors.Error(
						"Could not set roles list data",
						apperrors.WithCause(err),
					), nil)
				}
			}
		}()
		wg.Wait()
	})

	return w
}
