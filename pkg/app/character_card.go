package app

import (
	"context"
	"encoding/base64"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/esi"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

var ImagePlaceholder []byte

func init() {
	data := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="
	ImagePlaceholder = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	_, _ = base64.StdEncoding.Decode(ImagePlaceholder, []byte(data))
}

type CharacterCard struct {
	widget.BaseWidget

	deps   dependencies
	char   DataProxy[*repository.CharacterDBData]
	tags   *DataList[*repository.TagDBData]
	parent fyne.Window

	NameLabel         *widget.Label
	Portrait          *canvas.Image
	CorporationLabel  *widget.Label
	CorporationTicker *widget.Label
	CorporationIcon   *canvas.Image
	AllianceLabel     *widget.Label
	AllianceTicker    *widget.Label
	AllianceIcon      *canvas.Image
	RefreshButton     *widget.Button
	DeleteButton      *widget.Button

	MiniTags          *MiniTagSet[string, *CharacterMiniTag]
	MiniTagsContainer *container.Scroll
	miniTagLookup     map[int64]bool
	selectedTags      map[string]bool
}

type images struct {
	Portrait    *canvas.Image
	Corporation *canvas.Image
	Alliance    *canvas.Image
}

func getImagesForChar(logger logging.Logger, char *repository.CharacterDBData) (im images) {
	if char == nil {
		return images{
			Portrait:    canvas.NewImageFromResource(fyne.NewStaticResource("placeholder", ImagePlaceholder)),
			Corporation: canvas.NewImageFromResource(fyne.NewStaticResource("placeholder", ImagePlaceholder)),
			Alliance:    canvas.NewImageFromResource(fyne.NewStaticResource("placeholder", ImagePlaceholder)),
		}
	}

	portaitURI, err := storage.ParseURI(char.Character.Picture)
	if err != nil {
		level.Info(logger).Err("could not parse portait", err)
		im.Portrait = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder", ImagePlaceholder))
	} else {
		im.Portrait = canvas.NewImageFromURI(portaitURI)
	}

	corpIconURI, err := storage.ParseURI(char.Corporation.Picture)
	if err != nil {
		level.Info(logger).Err("could not parse corporation icon", err)
		im.Corporation = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder", ImagePlaceholder))
	} else {
		im.Corporation = canvas.NewImageFromURI(corpIconURI)
	}

	allianceIconURI, err := storage.ParseURI(char.Alliance.Picture.String)
	if err != nil {
		level.Info(logger).Err("could not parse alliance icon", err)
		im.Alliance = canvas.NewImageFromResource(fyne.NewStaticResource("placeholder", ImagePlaceholder))
	} else {
		im.Alliance = canvas.NewImageFromURI(allianceIconURI)
	}

	return im
}

func NewCharacterCard(deps dependencies, parent fyne.Window, dataChar DataProxy[*repository.CharacterDBData], tags *DataList[*repository.TagDBData], deleteFunc func(c *CharacterCard)) *CharacterCard {
	logger := logging.With(deps.Logger(), keys.Component, "CharacterCard.NewCharacterCard")

	char, err := dataChar.Get()
	if err != nil || char == nil {
		ShowError(logger, parent, AppError(
			"Could not find character data",
			WithCause(err),
		), nil)
		return nil
	}

	images := getImagesForChar(logger, char)

	corpTickText := fmt.Sprintf("[%s]", char.Corporation.Ticker)
	allyTickText := ""
	if char.Alliance.Ticker.String != "" {
		allyTickText = fmt.Sprintf("[%s]", char.Alliance.Ticker.String)
	}

	cc := &CharacterCard{
		deps:   deps,
		char:   dataChar,
		tags:   tags,
		parent: parent,

		NameLabel:         widget.NewLabel(char.Character.Name),
		Portrait:          images.Portrait,
		CorporationLabel:  widget.NewLabel(char.Corporation.Name),
		CorporationTicker: widget.NewLabel(corpTickText),
		CorporationIcon:   images.Corporation,
		AllianceLabel:     widget.NewLabel(char.Alliance.Name.String),
		AllianceTicker:    widget.NewLabel(allyTickText),
		AllianceIcon:      images.Alliance,
		// RefreshButton:     widget.NewButtonWithIcon("refresh", theme.ViewRefreshIcon(), nil),
		// DeleteButton:      widget.NewButtonWithIcon("delete", theme.DeleteIcon(), nil),
		RefreshButton: widget.NewButton("refresh", nil),
		DeleteButton:  widget.NewButton("delete", nil),

		MiniTags:      NewMiniTagSet[string, *CharacterMiniTag](),
		miniTagLookup: map[int64]bool{},
		selectedTags:  map[string]bool{},
	}
	cc.ExtendBaseWidget(cc)

	cc.MiniTagsContainer = container.NewVScroll(cc.MiniTags)

	cc.Portrait.FillMode = canvas.ImageFillStretch
	cc.Portrait.SetMinSize(fyne.Size{Height: 128, Width: 128})

	cc.CorporationIcon.FillMode = canvas.ImageFillStretch
	cc.CorporationIcon.SetMinSize(fyne.Size{Height: 64, Width: 64})

	cc.AllianceIcon.FillMode = canvas.ImageFillStretch
	cc.CorporationIcon.SetMinSize(fyne.Size{Height: 64, Width: 64})

	cc.RefreshButton.OnTapped = cc.refreshData
	cc.DeleteButton.OnTapped = cc.deleteCharacter(deleteFunc)
	cc.DeleteButton.Importance = widget.DangerImportance

	cc.refreshTags()

	dataChar.AddListener(binding.NewDataListener(cc.redraw))
	tags.AddListener(binding.NewDataListener(cc.refreshTags))

	return cc
}

// func (c *CharacterCard) Tapped(_ *fyne.PointEvent) {}

func (c *CharacterCard) CreateRenderer() fyne.WidgetRenderer {
	return c
}

func (c *CharacterCard) matchesSelectedTags() bool {
	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterCard.matchesSelectedTags")

	matchedTags := make(map[string]bool, c.MiniTags.Len())
	for _, mt := range c.MiniTags.items {
		tag, err := mt.tag.Get()
		if err != nil {
			ShowError(logger, c.parent, AppError(
				"Could not load tag data",
				WithCause(err),
			), nil)
		}
		matchedTags[tag.StrID()] = mt.isMatch
	}

	for k, on := range c.selectedTags {
		if !on {
			continue
		}

		if !matchedTags[k] {
			return false
		}
	}

	return true
}

func (c *CharacterCard) redraw() {
	defer c.Refresh()

	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterCard.redraw")

	char, err := c.char.Get()
	if err != nil || char == nil {
		ShowError(logger, c.parent, AppError(
			"Could not find character data",
			WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.CharacterID, char.Character.ID, keys.CharacterName, char.Character.Name)

	level.Debug(logger).Message("refreshing character card")

	if c.matchesSelectedTags() {
		if c.Hidden {
			c.Show()
		}
	} else {
		if !c.Hidden {
			c.Hide()
		}
	}

	c.NameLabel.Text = char.Character.Name
	c.CorporationLabel.Text = char.Corporation.Name
	c.AllianceLabel.Text = char.Alliance.Name.String

	c.CorporationTicker.Text = fmt.Sprintf("[%s]", char.Corporation.Ticker)
	c.AllianceTicker.Text = ""
	if char.Alliance.Ticker.String != "" {
		c.AllianceTicker.Text = fmt.Sprintf("[%s]", char.Alliance.Ticker.String)
	}

	images := getImagesForChar(logger, char)
	c.Portrait.Resource = images.Portrait.Resource
	c.CorporationIcon.Resource = images.Corporation.Resource
	c.AllianceIcon.Resource = images.Alliance.Resource
}

func (c *CharacterCard) refreshTags() {
	defer c.Refresh()

	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterCard.refreshTags")

	char, err := c.char.Get()
	if err != nil || char == nil {
		ShowError(logger, c.parent, AppError(
			"Could not find character data",
			WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.CharacterID, char.Character.ID, keys.CharacterName, char.Character.Name)

	tags, err := c.tags.Get()
	if err != nil {
		ShowError(logger, c.parent, AppError(
			"Could not find tags data",
			WithCause(err),
			WithInternalData(keys.CharacterID),
		), nil)
		return
	}

	for i, tag := range tags {
		_, ok := c.miniTagLookup[tag.Tag.ID]
		if ok { // old -- refresh triggered elsewise
			continue
		}

		c.miniTagLookup[tag.Tag.ID] = true
		mt := NewCharacterMiniTag(c.deps, c.parent, c.char, c.tags.Child(i))
		c.MiniTags.Add(mt)
		mt.Refresh()
	}
	c.MiniTags.Refresh()
}

func (c *CharacterCard) CharacterID() int64 {
	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterCard.CharacterID")

	char, err := c.char.Get()
	if err != nil {
		ShowError(logger, c.parent, AppError(
			"Could not find character data",
			WithCause(err),
		), nil)
		return 0
	}
	return char.Character.ID
}

func (c *CharacterCard) RefreshDataWith(char repository.CharacterDBData) error {
	return c.char.Set(&char)
}

func (c *CharacterCard) refreshData() {
	ctx := context.Background()

	c.RefreshButton.Disable()
	defer c.RefreshButton.Enable()

	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterCard.refreshData")

	char, err := c.char.Get()
	if err != nil {
		ShowError(logger, c.parent, AppError(
			"Could not find character data",
			WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.CharacterID, char.Character.ID, keys.CharacterName, char.Character.Name)

	dbTok, err := c.deps.AppRepo().GetTokenForCharacter(ctx, char.Character.ID, nil)
	if err != nil {
		ShowError(logger, c.parent, AppError(
			"Could not find character token",
			WithCause(err),
		), nil)
		return
	}

	tok := esi.TokenFromRepository(dbTok)

	newChar, err := RefreshCharacterData(ctx, c.deps, tok, char.Character.ID)
	if err != nil {
		ShowError(logger, c.parent, AppError(
			"Error refreshing character data",
			WithCause(err),
		), nil)
		return
	}

	if err := c.char.Set(&newChar); err != nil {
		ShowError(logger, c.parent, AppError(
			"Could not set character data",
			WithCause(err),
		), nil)
		return
	}
}

func (c *CharacterCard) deleteCharacter(callback func(*CharacterCard)) func() {
	return func() {
		ctx := context.Background()

		c.DeleteButton.Disable()
		defer c.DeleteButton.Enable()

		logger := logging.With(c.deps.Logger(), keys.Component, "CharacterCard.deleteCharacter")

		char, err := c.char.Get()
		if err != nil {
			ShowError(logger, c.parent, AppError(
				"Could not find character data",
				WithCause(err),
			), nil)
			return
		}

		logger = logging.With(logger, keys.CharacterID, char.Character.ID, keys.CharacterName, char.Character.Name)

		conf := dialog.NewConfirm("Delete Character?", fmt.Sprintf("Are you sure you want to delete the character '%s'?", char.Character.Name), func(ok bool) {
			if !ok {
				return
			}

			if err := c.deps.AppRepo().DeleteCharacter(ctx, c.CharacterID(), nil); err != nil {
				ShowError(logger, c.parent, AppError(
					"Unable to delete character",
					WithCause(err),
				), nil)
			} else {
				callback(c)
			}
		}, c.parent)
		conf.SetConfirmImportance(widget.DangerImportance)
		conf.Show()
	}
}

// The WidgetRenderer interface

func (c *CharacterCard) Destroy() {}

func (c *CharacterCard) Layout(sz fyne.Size) {
	fontSize := fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText)

	// Pictures
	c.Portrait.Move(fyne.Position{X: 0, Y: 0})
	c.Portrait.Resize(fyne.Size{Height: 128, Width: 128})

	c.AllianceIcon.Move(fyne.Position{X: 128 + theme.Padding(), Y: 96})
	c.AllianceIcon.Resize(fyne.Size{Width: 32, Height: 32})

	c.CorporationIcon.Move(fyne.Position{X: 160 + 2*theme.Padding(), Y: 96})
	c.CorporationIcon.Resize(fyne.Size{Width: 32, Height: 32})

	c.CorporationLabel.Hide()
	c.AllianceLabel.Hide()

	// Name and tickers
	c.NameLabel.Alignment = fyne.TextAlignLeading
	c.NameLabel.Wrapping = fyne.TextWrapOff
	c.NameLabel.TextStyle.Bold = true
	c.NameLabel.Move(fyne.Position{X: 128 + theme.Padding(), Y: 0})
	c.NameLabel.Refresh()
	namesz := fyne.MeasureText(c.NameLabel.Text, fontSize, c.NameLabel.TextStyle)

	c.AllianceTicker.Alignment = fyne.TextAlignLeading
	c.AllianceTicker.Wrapping = fyne.TextWrapOff
	c.AllianceTicker.TextStyle.Bold = true
	c.AllianceTicker.Move(fyne.Position{X: 128 + theme.Padding() + namesz.Width + theme.InnerPadding(), Y: 0}) // no right-side inner padding
	allysz := fyne.MeasureText(c.AllianceTicker.Text, fontSize, c.AllianceTicker.TextStyle)

	c.CorporationTicker.Alignment = fyne.TextAlignLeading
	c.CorporationTicker.Wrapping = fyne.TextWrapOff
	c.CorporationTicker.Move(fyne.Position{X: 128 + theme.Padding() + namesz.Width + theme.InnerPadding() + theme.Padding() + allysz.Width + theme.InnerPadding()/4, Y: 0}) // no right-side inner padding

	// Tags
	x := 128 + theme.Padding() + theme.InnerPadding()             // line up with the name
	y := namesz.Height + theme.InnerPadding() + theme.Padding()/2 // no extra padding
	c.MiniTagsContainer.Move(fyne.Position{X: x, Y: y})
	c.MiniTagsContainer.Resize(fyne.Size{Width: sz.Width - x, Height: sz.Height - y - 32 - theme.Padding()/2})
	c.MiniTags.Resize(fyne.Size{Width: sz.Width - x, Height: sz.Height - y - 32 - theme.Padding()/2})

	// Buttons
	refreshLabelSz := fyne.MeasureText(c.RefreshButton.Text, fontSize, c.NameLabel.TextStyle)
	c.RefreshButton.Resize(fyne.Size{Width: refreshLabelSz.Width + 2*theme.InnerPadding(), Height: refreshLabelSz.Height + theme.InnerPadding()})
	c.RefreshButton.Alignment = widget.ButtonAlignCenter
	rbsz := c.RefreshButton.Size()
	c.RefreshButton.Move(fyne.Position{X: sz.Width - rbsz.Width, Y: sz.Height - rbsz.Height})

	deleteLabelSz := fyne.MeasureText(c.DeleteButton.Text, fontSize, c.NameLabel.TextStyle)
	c.DeleteButton.Resize(fyne.Size{Width: deleteLabelSz.Width + 2*theme.InnerPadding(), Height: deleteLabelSz.Height + theme.InnerPadding()})
	c.DeleteButton.Alignment = widget.ButtonAlignCenter
	dbsz := c.DeleteButton.Size()
	c.DeleteButton.Move(fyne.Position{X: sz.Width - rbsz.Width - theme.Padding() - dbsz.Width, Y: sz.Height - dbsz.Height})
}

func (c *CharacterCard) MinSize() fyne.Size {
	return fyne.Size{
		Height: 128,
		Width:  256,
	}
}

func (c *CharacterCard) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{
		c.Portrait,
		c.NameLabel,
		c.CorporationIcon,
		c.CorporationLabel,
		c.CorporationTicker,
		c.AllianceIcon,
		c.AllianceLabel,
		c.AllianceTicker,
		c.RefreshButton,
		c.DeleteButton,
		c.MiniTagsContainer,
	}

	return objs
}

func (c *CharacterCard) Refresh() {
	for _, o := range c.Objects() {
		o.Refresh()
	}
}
