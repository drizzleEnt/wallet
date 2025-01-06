package service

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/drizzleent/wallet/models"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)
//0x140133C4cd251ef34DD884248f25C964dC75f0A6
//0x5295AFCE96E05C716d3C415236572DBAB9b5dA92i
const (
	dir = "./keystore"
)

type Service interface {
	SaveWallet(privatekey *ecdsa.PrivateKey, password string) error
	LoadWalletsFromKeystore() ([]models.KeystoreWallet, error)
	HasWallets() (bool, error)
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

func (b *service) LoadWalletsFromKeystore() ([]models.KeystoreWallet, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Failed open keystore folder %s", err.Error())
	}

	wallets := make([]models.KeystoreWallet, 0)
	for _, f := range entries {
		if f.IsDir() {
			continue
		}
		filePath := filepath.Join(dir, f.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("Failed open file %s", err.Error())
		}

		var wallet models.KeystoreWallet
		err = json.Unmarshal(data, &wallet)
		if err != nil {
			return nil, fmt.Errorf("Failed unmarshal file %s", err.Error())
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func (b *service) HasWallets() (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return false, nil
	}

	for _, f := range entries {
		if f.IsDir() {
			continue
		}
		filePath := filepath.Join(dir, f.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			return false, nil
		}

		var wallet models.KeystoreWallet
		err = json.Unmarshal(data, &wallet)
		if err == nil {
			return true, nil
		}
		log.Printf("ERROR: %s\n", err.Error())
	}
	return false, nil
}
