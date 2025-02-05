﻿package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/vitusb/wormhole-gui/v2/internal/assets"
)

const version = "v3.0.0-dev"
const commit = "3ec89d7e908f887786bb83dc5496ff610c2f613e"

var releaseURL = &url.URL{
	Scheme: "https",
	Host:   "github.com",
	// Path:   "/Jacalz/wormhole-gui/releases/tag/" + version,
    // Path:   "/vitusb/wormhole-gui/commit/" + commit,
    Path:   "/vitusb/wormhole-gui/",
}

type about struct {
	icon        *canvas.Image
	nameLabel   *widget.Label
    nameLabelVB *widget.Label
	spacerLabel *widget.Label
	hyperlink   *widget.Hyperlink
}

func newAbout() *about {
	return &about{}
}

func (a *about) buildUI() *fyne.Container {
	a.icon = canvas.NewImageFromResource(assets.AppIcon)
	a.icon.SetMinSize(fyne.NewSize(256, 256))

	a.nameLabel = newBoldLabel("Magic Wormhole Gui")
    a.nameLabelVB = newBoldLabel("Deutsche Version von \"Veit Berwig\" ...")
	a.spacerLabel = newBoldLabel("-")
	a.hyperlink = &widget.Hyperlink{Text: version, URL: releaseURL, TextStyle: fyne.TextStyle{Bold: true}}

	return container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), a.icon, layout.NewSpacer()),
		container.NewHBox(
			layout.NewSpacer(),
			a.nameLabel,
			a.spacerLabel,
			a.hyperlink,
            layout.NewSpacer(),
		),
        container.NewHBox(layout.NewSpacer(), a.nameLabelVB, layout.NewSpacer()),
		layout.NewSpacer(),
	)
}

func (a *about) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Über", Icon: theme.InfoIcon(), Content: a.buildUI()}
}
