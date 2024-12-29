package blockchain

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Blockchain interface {
	ImportFromPrivatekey(privatekeyHex string, password string) (string, *ecdsa.PrivateKey, error)
	ImportFromSeedPhrase(mnemonic string, password string) (string, *ecdsa.PrivateKey, error)
}

type blockchain struct {
}

func NewBlockchain() Blockchain {
	return &blockchain{}
}

func (b *blockchain) SendTransaction() {

}

func (b *blockchain) ImportFromSeedPhrase(mnemonic string, password string) (string, *ecdsa.PrivateKey, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, fmt.Errorf("mnemonic phrase not valid")
	}

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed create master key from seed")
	}

	childKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed create child key + 44")
	}
	childKey, err = masterKey.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed create child key + 60")
	}
	childKey, err = masterKey.NewChildKey(bip32.FirstHardenedChild)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed create child key + 0")
	}

	childKey, err = childKey.NewChildKey(0)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed create child key from child key")
	}

	childKey, err = childKey.NewChildKey(0)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed create child key from child key")
	}

	privateKey, err := crypto.ToECDSA(childKey.Key)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed create private key from child key")
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return address, privateKey, nil
}

func (b *blockchain) ImportFromPrivatekey(privatekeyHex string, password string) (string, *ecdsa.PrivateKey, error) {
	if privatekeyHex == "" {
		return "", nil, fmt.Errorf("Private key not set")
	}

	if strings.HasPrefix(privatekeyHex, "0x") {
		privatekeyHex = privatekeyHex[2:]
	}

	privateKeyBytes, err := hex.DecodeString(privatekeyHex)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed decode Private key to bytes")
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return "", nil, fmt.Errorf("Failed covert Private key to ECDSA")
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	log.Printf("wallet imported %s", address)

	return address, privateKey, nil
}
