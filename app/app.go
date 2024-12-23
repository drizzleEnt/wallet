package app

import (
	"github.com/drizzleent/wallet/blockchain"
	"github.com/drizzleent/wallet/service"
	"github.com/drizzleent/wallet/ui"
)

type App struct {
	bc      blockchain.Blockchain
	service service.Service
	ui      *ui.UI
}

func NewApp() *App {
	a := App{}
	return &a
}

func (a *App) Run() error {
	a.UI().RunApp()
	return nil
}

func (a *App) UI() *ui.UI {
	if a.ui == nil {
		a.ui = ui.NewUI()
	}
	return a.ui
}

func (a *App) Service() service.Service {
	if a.service == nil {
		a.service = service.NewService()
	}
	return a.service
}

func (a *App) Blockchain() blockchain.Blockchain {
	if a.bc == nil {
		a.bc = blockchain.NewBlockchain()
	}
	return a.bc
}
