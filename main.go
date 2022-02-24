package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/vitusb/wormhole-gui/v2/internal/assets"
	"github.com/vitusb/wormhole-gui/v2/internal/ui"
)

func main() {
	a := app.NewWithID("io.github.jacalz.wormhole_gui")
	a.SetIcon(assets.AppIcon)
	w := a.NewWindow("Magic Wormhole Gui")

	w.SetContent(ui.Create(a, w))
	w.Resize(fyne.NewSize(700, 400))
	w.SetMaster()
	w.ShowAndRun()
}
