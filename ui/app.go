package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
}

func NewUI() *UI {
	return &UI{}
}

func (a *UI) RunApp() {
	myApp := app.New()
	myWindow := myApp.NewWindow("I-WALLET")
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(800, 500))

	label := widget.NewLabel("Welcome to Crypto Wallet")
	createBtn := widget.NewButton("Import Wallet", func() {

		label.SetText("wallet imported")
	})

	content := container.NewVBox(label, createBtn)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
