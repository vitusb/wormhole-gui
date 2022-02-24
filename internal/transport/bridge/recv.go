﻿package bridge

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/vitusb/wormhole-gui/v2/internal/transport"
	"github.com/vitusb/wormhole-gui/v2/internal/util"
)

// RecvItem is the item that is being received
type RecvItem struct {
	URI      fyne.URI
	Progress *util.ProgressBar
	Name     string
}

// RecvList is a list of progress bars that track send progress.
type RecvList struct {
	widget.List

	client *transport.Client

	Items []*RecvItem
	lock  sync.RWMutex

	window fyne.Window
}

// Length returns the length of the data.
func (p *RecvList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *RecvList) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{},
		&widget.FileIcon{URI: nil},
		&widget.Label{Text: "Warte auf Daten ...", Wrapping: fyne.TextTruncate},
		util.NewProgressBar(),
	)
}

// UpdateItem updates the data in the list.
func (p *RecvList) UpdateItem(i int, item fyne.CanvasObject) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	container := item.(*fyne.Container)
	container.Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	container.Objects[1].(*widget.Label).SetText(p.Items[i].Name)
	p.Items[i].Progress = container.Objects[2].(*util.ProgressBar)
}

// OnSelected currently just makes sure that we don't persist selection.
func (p *RecvList) OnSelected(i int) {
	p.Unselect(i)
}

// NewRecvItem creates a new send item and adds it to the items.
func (p *RecvList) NewRecvItem() *RecvItem {
	p.lock.Lock()
	defer p.lock.Unlock()

	item := &RecvItem{Name: "Warte auf Daten ..."}
	p.Items = append(p.Items, item)
	return item
}

// NewReceive adds data about a new send to the list and then returns the channel to update the code.
func (p *RecvList) NewReceive(code string) {
	item := p.NewRecvItem()
	p.Refresh()

	path := make(chan string)

	go func() {
		name := <-path
		item.URI = storage.NewFileURI(name)
		if name != "text" {
			item.Name = item.URI.Name()
		} else {
			item.Name = "Text-Übertragung"
		}

		close(path)
		p.Refresh()
	}()

	go func(code string) {
		if err := p.client.NewReceive(code, path, item.Progress); err != nil {
			p.client.ShowNotification("Empfang fehlgeschlagen", "Ein Fehler ist beim Empfang der Daten aufgetreten.")
			item.Progress.Failed()
			dialog.ShowError(err, p.window)
		} else {
			p.client.ShowNotification("Empfang abgeschlossen", "Die Daten wurden erfolgreich empfangen.")
		}

		p.Refresh()
	}(code)
}

// NewRecvList greates a list of progress bars.
func NewRecvList(window fyne.Window, client *transport.Client) *RecvList {
	p := &RecvList{client: client, window: window}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelected = p.OnSelected
	p.ExtendBaseWidget(p)

	return p
}
