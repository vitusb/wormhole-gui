﻿package ui

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

type send struct {
	contentPicker dialog.Dialog

	fileChoice      *widget.Button
	directoryChoice *widget.Button
	textChoice      *widget.Button
	codeChoice      *widget.Check

	fileDialog      *dialog.FileDialog
	directoryDialog *dialog.FileDialog

	contentToSend *widget.Button
	sendList      *bridge.SendList

	client *transport.Client
	window fyne.Window
	app    fyne.App
}

func newSend(a fyne.App, w fyne.Window, c *transport.Client) *send {
	return &send{app: a, window: w, client: c}
}

func (s *send) onFileSend() {
	s.contentPicker.Hide()
	s.fileDialog.Resize(util.WindowSizeToDialog(s.window.Canvas().Size()))
	s.fileDialog.Show()
}

func (s *send) onDirSend() {
	s.contentPicker.Hide()
	s.fileDialog.Resize(util.WindowSizeToDialog(s.window.Canvas().Size()))
	s.directoryDialog.Show()
}

func (s *send) onTextSend() {
	s.contentPicker.Hide()
	s.sendList.SendText()
}

func (s *send) onCustomCode(enabled bool) {
	s.client.CustomCode = enabled
}

func (s *send) buildUI() *fyne.Container {
	s.fileChoice = &widget.Button{Text: "Datei", Icon: theme.FileIcon(), OnTapped: s.onFileSend}
	s.directoryChoice = &widget.Button{Text: "Verzeichnis", Icon: theme.FolderOpenIcon(), OnTapped: s.onDirSend}
	s.textChoice = &widget.Button{Text: "Text", Icon: theme.DocumentCreateIcon(), OnTapped: s.onTextSend}
	s.codeChoice = &widget.Check{Text: "Eigenen Code verwenden", OnChanged: s.onCustomCode}

	choiceContent := container.NewGridWithColumns(1, s.fileChoice, s.directoryChoice, s.textChoice, s.codeChoice)
	s.contentPicker = dialog.NewCustom("Was möchten Sie übertragen ?", "Abbrechen", choiceContent, s.window)

	s.sendList = bridge.NewSendList(s.window, s.client)
	s.contentToSend = &widget.Button{Text: "Inhalt zum Senden hinzufügen", Icon: theme.ContentAddIcon(), OnTapped: s.contentPicker.Show}

	s.fileDialog = dialog.NewFileOpen(s.sendList.OnFileSelect, s.window)
	s.directoryDialog = dialog.NewFolderOpen(s.sendList.OnDirSelect, s.window)

	box := container.NewVBox(s.contentToSend, &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, s.sendList)
}

func (s *send) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Senden", Icon: theme.MailSendIcon(), Content: s.buildUI()}
}
