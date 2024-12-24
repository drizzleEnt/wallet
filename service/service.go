package service

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

type Service interface {
	SaveWallet(privatekey string, password string) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (b *service) SaveWallet(privatekeyHex string, password string) error {
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	privateKeyBytes, err := hex.DecodeString(privatekeyHex)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return fmt.Errorf("Failed decode Private key to bytes")
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return fmt.Errorf("Failed covert Private key to ECDSA")
	}

	ac, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return fmt.Errorf("Failed create ks account")
	}
	log.Printf("wallet saved %s", ac.Address.Hex())
	return nil
}
