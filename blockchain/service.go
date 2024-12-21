package blockchain

type Blockchain interface {
	//CreateWallet()
	//SendTransaction()
	//GetBalance()
}

type blockchain struct {
}

func NewBlockchain() Blockchain {
	return &blockchain{}
}

func (b *blockchain) SendTransaction() {

}
