package blockchain

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/drizzleent/wallet/config"
	"github.com/drizzleent/wallet/models"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Blockchain interface {
	ImportFromPrivatekey(privatekeyHex string, password string) (string, *ecdsa.PrivateKey, error)
	ImportFromSeedPhrase(mnemonic string, password string) (string, *ecdsa.PrivateKey, error)
	CreateWallet() (string, *ecdsa.PrivateKey, error)
	GetEtherBalance(address string) (string, error)
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

func (b *blockchain) CreateWallet() (string, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return "", nil, fmt.Errorf("Failed create new wallet")
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	return address, privateKey, nil
}

func (b *blockchain) GetEtherBalance(address string) (string, error) {
	query := url.Values{}
	query.Add(`action`, `balance`)
	query.Add(`module`, `account`)
	query.Add(`address`, "0x"+address)
	query.Add(`tag`, `latest`)
	query.Add(`apikey`, config.API)

	u := url.URL{
		Scheme: `https`,
		Host:   config.EtherscanApi,
		Path:   `api`,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return "", fmt.Errorf("Failed get balance")
	}

	req.URL.RawQuery = query.Encode()
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return "", fmt.Errorf("Failed get balance")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return "", fmt.Errorf("Failed get balance")
	}

	var balance models.Balance

	err = json.Unmarshal(body, &balance)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return "", fmt.Errorf("Failed get balance")
	}

	weiFloat, err := strconv.ParseFloat(balance.Result, 64)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return "", fmt.Errorf("Failed get balance")
	}

	eth := b.ConvertAmountFromWei(big.NewInt(int64(weiFloat)))

	fmt.Printf("eth: %v\n", eth)

	return eth.String(), nil
}

func (bs *blockchain) ConvertAmountIntoWei(amount float64) *big.Int {
	oneEthInWei := new(big.Float).SetFloat64(math.Pow10(18))
	result, _ := new(big.Float).Mul(oneEthInWei, new(big.Float).SetFloat64(amount)).Int(new(big.Int))

	return result
}

func (bs *blockchain) ConvertAmountFromWei(amount *big.Int) *big.Float {
	oneEthInWei := new(big.Float).SetFloat64(math.Pow10(-18))
	result := new(big.Float).Mul(oneEthInWei, new(big.Float).SetInt(amount))

	return result
}
