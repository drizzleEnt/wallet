package app

import (
	"github.com/drizzleent/wallet/blockchain"
	"github.com/drizzleent/wallet/blockchain/evm"
	"github.com/drizzleent/wallet/config"
	"github.com/drizzleent/wallet/service"
	"github.com/drizzleent/wallet/ui"
)

type serviceProvider struct {
	evmCfg config.EvmConfig

	bc      blockchain.Blockchain
	service service.Service
	ui      *ui.UI
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (a *serviceProvider) UI() *ui.UI {
	if a.ui == nil {
		a.ui = ui.NewUI(a.Blockchain(), a.Service())
	}

	return a.ui
}

func (a *serviceProvider) Service() service.Service {
	if a.service == nil {
		a.service = service.NewService()
	}

	return a.service
}

func (a *serviceProvider) EvnConfig() config.EvmConfig {
	if a.evmCfg == nil {

	}

	return a.evmCfg
}

func (a *serviceProvider) Blockchain() blockchain.Blockchain {
	if a.bc == nil {
		a.bc = evm.NewBlockchain()
	}
	return a.bc
}
