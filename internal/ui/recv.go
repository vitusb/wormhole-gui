package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/vitusb/wormhole-gui/v2/internal/transport"
	"github.com/vitusb/wormhole-gui/v2/internal/transport/bridge"
	"github.com/vitusb/wormhole-gui/v2/internal/util"
)

type recv struct {
	codeEntry  *widget.Entry
	codeButton *widget.Button

	recvList *bridge.RecvList

	client *transport.Client
	window fyne.Window
	app    fyne.App
}

func newRecv(a fyne.App, w fyne.Window, c *transport.Client) *recv {
	return &recv{app: a, window: w, client: c}
}

func (r *recv) onRecv() {
	if err := r.codeEntry.Validate(); err != nil || r.codeEntry.Text == "" {
		dialog.ShowInformation("Ungültiger Code", "Der Code ist ungültig. Bitte erneut versuchen ...", r.window)
		return
	}

	r.recvList.NewReceive(r.codeEntry.Text)
	r.codeEntry.SetText("")
}

func (r *recv) buildUI() *fyne.Container {
	r.codeEntry = &widget.Entry{PlaceHolder: "Code eingeben", Wrapping: fyne.TextTruncate, OnSubmitted: func(_ string) { r.onRecv() },
		Validator: util.CodeValidator,
	}

	r.codeButton = &widget.Button{Text: "Herunterladen", Icon: theme.DownloadIcon(), OnTapped: r.onRecv}

	r.recvList = bridge.NewRecvList(r.window, r.client)

	box := container.NewVBox(container.NewGridWithColumns(2, r.codeEntry, r.codeButton), &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, r.recvList)
}

func (r *recv) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Empfangen", Icon: theme.DownloadIcon(), Content: r.buildUI()}
}
