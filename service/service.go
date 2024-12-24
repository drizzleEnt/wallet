package service

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type Service interface {
	SaveWallet(privatekey *ecdsa.PrivateKey, password string) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (b *service) SaveWallet(privateKey *ecdsa.PrivateKey, password string) error {
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	ac, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return fmt.Errorf("Failed create ks account")
	}
	log.Printf("wallet saved %s", ac.Address.Hex())
	return nil
}
