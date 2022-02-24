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

// SendItem is the item that is being sent.
type SendItem struct {
	URI      fyne.URI
	Progress *util.ProgressBar
	Code     string
	Name     string
}

// SendList is a list of progress bars that track send progress.
type SendList struct {
	widget.List

	client *transport.Client

	Items []*SendItem
	lock  sync.RWMutex

	window fyne.Window
}

// Length returns the length of the data.
func (p *SendList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *SendList) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{},
		&widget.FileIcon{URI: nil},
		&widget.Label{Text: "Warte auf Daten ...", Wrapping: fyne.TextTruncate},
		newCodeDisplay(p.window),
		util.NewProgressBar(),
	)
}

// UpdateItem updates the data in the list.
func (p *SendList) UpdateItem(i int, item fyne.CanvasObject) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	container := item.(*fyne.Container)
	container.Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	container.Objects[1].(*widget.Label).SetText(p.Items[i].Name)
	container.Objects[2].(*fyne.Container).Objects[0].(*codeDisplay).SetText(p.Items[i].Code)
	p.Items[i].Progress = container.Objects[3].(*util.ProgressBar)
}

// OnSelected currently just makes sure that we don't persist selection.
func (p *SendList) OnSelected(i int) {
	p.Unselect(i)
}

// NewSendItem adds data about a new send to the list and then returns the item.
func (p *SendList) NewSendItem(name string, uri fyne.URI) *SendItem {
	p.lock.Lock()
	defer p.lock.Unlock()

	item := &SendItem{Name: name, Code: "Warte auf Code ...", URI: uri}
	p.Items = append(p.Items, item)
	return item
}

// OnFileSelect is intended to be passed as callback to a FileOpen dialog.
func (p *SendList) OnFileSelect(file fyne.URIReadCloser, err error) {
	if err != nil {
		fyne.LogError("Error on selecting file to send", err)
		dialog.ShowError(err, p.window)
		return
	} else if file == nil {
		return
	}

	item := p.NewSendItem(file.URI().Name(), file.URI())
	p.Refresh()

	go func() {
		// We want to catch close errors for security reasons.
		defer func() {
			if err = file.Close(); err != nil {
				item.Progress.Failed()
				fyne.LogError("Error on closing file", err)
			}
		}()

		code, result, err := p.client.NewFileSend(file, item.Progress.WithProgress(), p.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending file", err)
			item.Progress.Failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending file", res.Error)
			item.Progress.Failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("Versenden der Datei fehlgeschlagen", "Ein Fehler ist beim Versenden aufgetreten.")
		} else if res.OK {
			p.client.ShowNotification("Versenden der Datei abgeschlossen", "Die Daten wurden erfolgreich übertragen.")
		}
	}()
}

// OnDirSelect is intended to be passed as callback to a FolderOpen dialog.
func (p *SendList) OnDirSelect(dir fyne.ListableURI, err error) {
	if err != nil {
		fyne.LogError("Error on selecting dir to send", err)
		dialog.ShowError(err, p.window)
		return
	} else if dir == nil {
		return
	}

	item := p.NewSendItem(dir.Name(), dir)
	p.Refresh()

	go func() {
		code, result, err := p.client.NewDirSend(dir, item.Progress.WithProgress(), p.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending directory", err)
			item.Progress.Failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending directory", res.Error)
			item.Progress.Failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("Versenden des Verzeichnisses fehlgeschlagen", "Ein Fehler ist beim Versenden aufgetreten.")
		} else if res.OK {
			p.client.ShowNotification("Versenden des Verzeichnisses abgeschlossen", "Die Daten wurden erfolgreich übertragen.")
		}
	}()
}

// SendText sends new text.
func (p *SendList) SendText() {
	// The file URI is a hack to get the correct icon.
	item := &SendItem{Name: "Text-Übertragung", Code: "Warte auf Code ...", URI: storage.NewFileURI("text")}

	go func() {
		text := <-p.client.ShowTextSendWindow()
		if text == "" {
			return
		}

		p.Items = append(p.Items, item)
		p.Refresh()

		code, result, err := p.client.NewTextSend(text, item.Progress.WithProgress(), p.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending text", err)
			item.Progress.Failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending text", res.Error)
			item.Progress.Failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("Versenden des Textes fehlgeschlagen", "Ein Fehler ist beim Versenden aufgetreten.")
		} else if res.OK && p.client.Notifications {
			p.client.ShowNotification("Versenden des Textes abgeschlossen", "Die Daten wurden erfolgreich übertragen.")
		}
	}()
}

// getCustomCode returns "" if the user has custom codes disabled.
// Otherwise, it will ask the user for a code.
func (p *SendList) getCustomCode() string {
	if !p.client.CustomCode {
		return ""
	}

	code := make(chan string)
	codeEntry := &widget.Entry{
		PlaceHolder: "123-insecure-example-code",
		Wrapping:    fyne.TextTruncate,
		Validator:   util.CodeValidator,
	}

	form := dialog.NewForm("Create custom code", "Confirm", "Cancel", []*widget.FormItem{
		{Text: "Code", Widget: codeEntry}, // TODO: Add HintText when FormDialog supports it.
	}, func(submitted bool) {
		if !submitted {
			code <- ""
		} else {
			code <- codeEntry.Text
		}

		close(code)
	}, p.window)
	form.Resize(fyne.Size{Width: p.window.Canvas().Size().Width * 0.8, Height: 0})
	form.Show()

	return <-code
}

// NewSendList greates a list of progress bars.
func NewSendList(window fyne.Window, client *transport.Client) *SendList {
	p := &SendList{client: client, window: window}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelected = p.OnSelected
	p.ExtendBaseWidget(p)

	return p
}
