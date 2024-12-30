package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
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

	hasWallets, err := a.srv.HasWallets()
	if err != nil {
		dialog.ShowError(err, myWindow)
		return
	}

	if hasWallets {
		a.showMainMenu(myWindow)
		//a.showWalletList(myWindow)
	} else {
		a.showStartWall(myWindow)
	}

	myWindow.ShowAndRun()
}

func (a *UI) showStartWall(w fyne.Window) {
	label := widget.NewLabelWithStyle(
		"Welcome to Crypto Wallet",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)
	//label := widget.NewLabel("Welcome to Crypto Wallet")
	importBtn := widget.NewButton("Import Wallet", func() {
		a.importWallet(w)
	})
	createBtn := widget.NewButton("Create Wallet", func() {
		a.createWallet(w)
	})
	content := container.NewVBox(label, importBtn, createBtn)
	w.SetContent(content)
}

func (a *UI) importWallet(w fyne.Window) {
	importPrivatekeyBtn := widget.NewButton("Import Private key", func() {
		a.showImportPrivatekey(w)
	})
	importSeedPhraseBtn := widget.NewButton("Import Seed phrase", func() {
		a.showImportSeedPhrase(w)
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle(
			"Wallet import",
			fyne.TextAlignCenter,
			fyne.TextStyle{},
		),
		importPrivatekeyBtn,
		importSeedPhraseBtn,
	)

	w.SetContent(content)
}

func (a *UI) createWallet(w fyne.Window) {

	dialog.ShowConfirm("Create New Wallet", "Create new wallet?", func(b bool) {
		if b {
			address, privateKey, err := a.bc.CreateWallet()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			fmt.Printf("address: %v\n", address)
			fmt.Printf("privateKey: %v\n", privateKey)
			a.showSaveWithPassword(w, func(password string) {
				err := a.srv.SaveWallet(privateKey, password)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				dialog.ShowInformation("Success", "Wallet created and saved successfully!", w)
				a.showMainMenu(w)
			})
		}
	}, w)
}

func (a *UI) showImportPrivatekey(w fyne.Window) {
	privatekeyEntry := widget.NewEntry()
	privatekeyEntry.SetPlaceHolder("Enter private key")

	dialog.ShowCustomConfirm("Import Private Key", "Confirm", "Cancel", privatekeyEntry, func(b bool) {
		if b {
			privatekeyStr := privatekeyEntry.Text
			if privatekeyStr == "" {
				dialog.ShowError(fmt.Errorf("Private key cannot be empty"), w)
				return
			}

			address, privatekey, err := a.bc.ImportFromPrivatekey(privatekeyStr, "")
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed import private key"), w)
				return
			}
			fmt.Printf("address: %v\n", address)

			a.showSaveWithPassword(w, func(password string) {
				err := a.srv.SaveWallet(privatekey, password)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed save wallet"), w)
					return
				}
				dialog.ShowInformation("Success", "Wallet imported and saved successfully!", w)
				a.showMainMenu(w)
			})
		}
	}, w)
}

func (a *UI) showImportSeedPhrase(w fyne.Window) {
	seedEntry := widget.NewEntry()
	seedEntry.SetPlaceHolder("Enter seed phrase")

	dialog.ShowCustomConfirm("Import Seed Phrase", "Confirm", "Cancel", seedEntry, func(b bool) {
		if b {
			seedStr := seedEntry.Text
			if seedStr == "" {
				dialog.ShowError(fmt.Errorf("Seed Phrase cannot be empty"), w)
				return
			}

			address, privatekey, err := a.bc.ImportFromSeedPhrase(seedStr, "")
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed import private key"), w)
				return
			}
			fmt.Printf("address: %v\n", address)

			a.showSaveWithPassword(w, func(password string) {
				err := a.srv.SaveWallet(privatekey, password)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed save wallet"), w)
					return
				}
				dialog.ShowInformation("Success", "Wallet imported and saved successfully!", w)
				a.showMainMenu(w)
			})
		}
	}, w)
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

func (a *UI) showWalletList(w fyne.Window) {
	wallets, err := a.srv.LoadWalletsFromKeystore()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	list := widget.NewList(
		func() int { return len(wallets) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(wallets[lii].Address)
		},
	)
	btns := container.NewCenter(container.NewHBox(
		widget.NewButton("Refresh", func() {
			a.showWalletList(w)
		}),
		widget.NewButton("Back", func() {

		})))

	content := container.NewBorder(
		nil,
		btns,
		nil,
		nil,
		list,
	)

	w.SetContent(content)
}

func (a *UI) showSettingsMenu(w fyne.Window) {
	label := widget.NewLabelWithStyle(
		"Settings",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)
	backBtn := widget.NewButtonWithIcon("Back", theme.CancelIcon(), func() {
		a.showMainMenu(w)
	})
	importBtn := widget.NewButton("Import Wallet", func() {
		a.importWallet(w)
	})
	createBtn := widget.NewButton("Create Wallet", func() {
		a.createWallet(w)
	})
	walletListBtn := widget.NewButton("Switch Wallet", func() {
		a.showWalletList(w)
	})

	center := container.NewVBox(
		importBtn,
		createBtn,
		walletListBtn,
	)
	content := container.NewBorder(
		label,
		backBtn,
		nil,
		nil,
		center,
	)
	w.SetContent(content)
}

func (a *UI) showMainMenu(w fyne.Window) {
	wallets, err := a.srv.LoadWalletsFromKeystore()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	totalBalance := "$1"

	tokenList := []string{
		"ARB: 50",
		"ETH: 1",
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			a.showSettingsMenu(w)
		}),
	)

	addressLabel := widget.NewLabelWithStyle(
		"Wallet: "+wallets[0].Address,
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	balanceLabel := widget.NewLabelWithStyle(
		totalBalance,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	tokenListWidget := widget.NewList(
		func() int { return len(tokenList) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(tokenList[lii])
		},
	)

	bottomSection := container.NewVBox(
		toolbar,
	)

	topSection := container.NewVBox(
		//toolbar,
		addressLabel,
		balanceLabel,
	)

	content := container.NewBorder(
		topSection,
		bottomSection,
		nil,
		nil,
		tokenListWidget,
	)
	w.SetContent(content)
}
