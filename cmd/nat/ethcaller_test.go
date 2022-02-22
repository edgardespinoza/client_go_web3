package nat

import (
	"encoding/hex"
	"log"
	"math/big"
	"testing"

	"edu.id/sdk/cmd/lib"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	HOST = "https://data-seed-prebsc-1-s1.binance.org:8545/"
)

func TestGetVersion(t *testing.T) {

	eth := NewEthBasicCaller(HOST, nil)
	defer eth.Close()

	value := eth.GetVersion()

	log.Println("version:", value)

}

func TestGetBlockNumber(t *testing.T) {

	eth := NewEthBasicCaller(HOST, nil)
	defer eth.Close()

	value := eth.GetBlockNumber()

	log.Println("price:", value)

}

func TestGetNonce(t *testing.T) {

	eth := NewEthBasicCaller(HOST, nil)
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered", r)
		}
	}()
	defer eth.Close()

	pk := "cd4d53fe14ef4943887b1f4be2bd175f9979ce3713bf96a9c32a9245fbdc9056"

	fromAddress := lib.PrivateToAddress(pk)

	log.Println("Send", fromAddress.Hex())

	nonce := eth.GetNonce(fromAddress)

	log.Println("nonce:", nonce)

}

func TestGetGasPrice(t *testing.T) {

	eth := NewEthBasicCaller(HOST, nil)
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered", r)
		}
	}()
	defer eth.Close()

	gasPrices := eth.GetGasPrice()

	log.Println("price:", gasPrices)

}

func TestEstimateGas(t *testing.T) {

	eth := NewEthBasicCaller(HOST, nil)
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered", r)
		}
	}()
	defer eth.Close()

	pk := "cd4d53fe14ef4943887b1f4be2bd175f9979ce3713bf96a9c32a9245fbdc9056"
	addressErc20 := common.HexToAddress("0xF89Ee199e4ed0500ce8cE245d59542d693A23AfC")

	to := common.HexToAddress("0xB3cedeebaea9F8DeA537cc8e6136B6338104FE0f")

	fromAddress := lib.PrivateToAddress(pk)

	log.Println("from", fromAddress.Hex())
	amount := big.NewInt(1)

	api, _ := BuildApi()
	data, _ := TransferByteCode(api, &to, amount)
	gasLimit, _ := eth.GetEstimateGas(fromAddress, &addressErc20, data)

	log.Println("gaslimit:", gasLimit)

}

func TestGetAddresstoPub(t *testing.T) {
	pub := "358bbe04823acb656d1635e6e89311c5202de5c738b24b7be06301babbe7a3b72bc3e58abe30a7510ea253157f852ea8603d767512668f2ac06f50d9c29e760a"

	address := lib.PublicToAddress(pub)

	log.Println("address:", address.Hex())

}

func TestSendTransaction(t *testing.T) {

	pkSenderAlice := "cd4d53fe14ef4943887b1f4be2bd175f9979ce3713bf96a9c32a9245fbdc9056"
	contractAddressErc20 := common.HexToAddress("0xF89Ee199e4ed0500ce8cE245d59542d693A23AfC")
	addressBob := common.HexToAddress("0xB3cedeebaea9F8DeA537cc8e6136B6338104FE0f")
	tokens := 1

	eth := NewEthBasicCaller(HOST, nil)

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered", r)
		}
	}()
	defer eth.Close()

	prvSender, err := crypto.HexToECDSA(pkSenderAlice)
	if err != nil {
		log.Fatal(err)
	}
	addressAlice := lib.PrivateToAddress(pkSenderAlice)

	decimal := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	amount := new(big.Int).Mul(decimal, big.NewInt(int64(tokens)))

	log.Println("total:", amount)

	api, err := BuildApi()
	if err != nil {
		log.Fatal(err)
	}
	data, err := TransferByteCode(api, &addressBob, amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("bytecode:", hex.EncodeToString(*data))
	gasLimit, err := eth.GetEstimateGas(addressAlice, &contractAddressErc20, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gaslimit:", gasLimit)

	gasPrice := eth.GetGasPrice()

	nonce := eth.GetNonce(addressAlice)

	tx, err := eth.GetTransaction(gasPrice, nonce, gasLimit, &contractAddressErc20, data)
	if err != nil {
		log.Fatal(err)
	}

	txSigner := eth.SignTransaction(tx, prvSender)

	hashTx := eth.Send(txSigner)

	log.Println("Hash tx:", hashTx)
}
