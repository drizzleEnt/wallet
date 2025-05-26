package blockchain

import (
	"crypto/ecdsa"
)

type Blockchain interface {
	ImportFromPrivatekey(privatekeyHex string, password string) (string, *ecdsa.PrivateKey, error)
	ImportFromSeedPhrase(mnemonic string, password string) (string, *ecdsa.PrivateKey, error)
	CreateWallet() (string, *ecdsa.PrivateKey, error)
	GetBalance(address string) (string, error)
}
