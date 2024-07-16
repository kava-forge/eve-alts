package app

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/keys"
)

func NewMainWindow(deps dependencies, a fyne.App) fyne.Window {
	logger := logging.With(deps.Logger(), keys.Component, "MainWindow")

	w := a.NewWindow("EVE Alts")
	w.CenterOnScreen()
	w.Resize(fyne.Size{Width: 1024, Height: 768})

	tagsTab, tags := NewTagsTab(deps, w)
	charsTab, chars := NewCharactersTab(deps, w, tags)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Characters", theme.AccountIcon(), charsTab),
		container.NewTabItemWithIcon("Tags", theme.ListIcon(), tagsTab),
	)

	w.SetContent(tabs)

	a.Lifecycle().SetOnStarted(func() {
		ctx := context.Background()
		knownChars, err := deps.AppRepo().GetAllCharacters(ctx, nil)
		if err != nil && !errors.Is(err, database.ErrNoRows) {
			ShowError(logger, w, AppError(
				"Could not load character data",
				WithCause(err),
			), nil)
		} else {
			level.Debug(logger).Message("initial characters loaded")
			if err := chars.Set(knownChars); err != nil {
				ShowError(logger, w, AppError(
					"Could not set character list data",
					WithCause(err),
				), nil)
			}
		}

		knownTags, err := deps.AppRepo().GetAllTags(ctx, nil)
		if err != nil && !errors.Is(err, database.ErrNoRows) {
			ShowError(logger, w, AppError(
				"Could not load tags data",
				WithCause(err),
			), nil)
		} else {
			level.Debug(logger).Message("initial tags loaded")
			if err := tags.Set(knownTags); err != nil {
				ShowError(logger, w, AppError(
					"Could not set tags list data",
					WithCause(err),
				), nil)
			}
		}
	})

	return w
}
