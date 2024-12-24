package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/drizzleent/wallet/blockchain"
	"github.com/drizzleent/wallet/service"
)

type UI struct {
	bc  blockchain.Blockchain
	srv service.Service
}

func NewUI(b blockchain.Blockchain, s service.Service) *UI {
	return &UI{
		bc:  b,
		srv: s,
	}
}

func (a *UI) RunApp() {
	myApp := app.New()
	myWindow := myApp.NewWindow("I-WALLET")
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(800, 500))

	label := widget.NewLabel("Welcome to Crypto Wallet")
	createBtn := widget.NewButton("Import Wallet", func() {
		a.ImportWallet(myWindow)
	})

	content := container.NewVBox(label, createBtn)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func (a *UI) ImportWallet(w fyne.Window) {
	importPrivatekeyBtn := widget.NewButton("Import Private key", func() {
		a.showImportPrivatekey(w)
	})
	importSeedPhraseBtn := widget.NewButton("Import Seed phrase", func() {
		a.showImportSeedPhrase(w)
	})

	content := container.NewVBox(
		widget.NewLabel("Wallet import"),
		importPrivatekeyBtn,
		importSeedPhraseBtn,
	)

	w.SetContent(content)
}

func (a *UI) showImportPrivatekey(w fyne.Window) {
	privatekeyEntry := widget.NewEntry()
	privatekeyEntry.SetPlaceHolder("Enter private key")

	dialog.ShowCustomConfirm("Import Private Key", "Confirm", "Cancel", privatekeyEntry, func(b bool) {
		if b {
			privatekey := privatekeyEntry.Text
			if privatekey == "" {
				dialog.ShowError(fmt.Errorf("Private key cannot be empty"), w)
				return
			}
			a.showSaveWithPassword(w, func(password string) {
				err := a.srv.SaveWallet(privatekey, password)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed save wallet"), w)
					return
				}
				dialog.ShowInformation("Success", "Wallet imported and saved successfully!", w)
			})
		}
	}, w)
}

func (a *UI) showImportSeedPhrase(w fyne.Window) {

}

func (a *UI) showSaveWithPassword(w fyne.Window, onSave func(password string)) {
	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Enter password")

	dialog.ShowCustomConfirm("Set Password", "Save", "Cancel", passwordEntry, func(b bool) {
		if b {
			password := passwordEntry.Text
			if password == "" {
				dialog.ShowError(fmt.Errorf("password cannot be empty"), w)
				return
			}
			onSave(password)
		}
	}, w)
}
