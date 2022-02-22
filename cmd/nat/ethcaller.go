package nat

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthCaller interface {
	Init()
	GetVersion() string
	GetGasPrice() int64
	GetNonce(address string) int64
	GetBlockNumber() int64
	Close()
}

type EthBasicCaller struct {
	host   string
	client *ethclient.Client
}

func NewEthBasicCaller(node string, headers map[string]string) *EthBasicCaller {
	e := EthBasicCaller{host: node}

	rpcClient, err := rpc.Dial(e.host)

	if err != nil {
		log.Println("Error init", err)
		log.Fatal(err)
	}

	for k, v := range headers {
		//rpcClient.SetHeader("Authorization", "Bearer "+token)
		rpcClient.SetHeader(k, v)
	}

	e.client = ethclient.NewClient(rpcClient)
	return &e
}

func (e *EthBasicCaller) GetVersion() string {
	version, err := e.client.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}
	log.Println("version: ", version)
	return fmt.Sprintf("0x%x", version)
}

func (e *EthBasicCaller) GetBlockNumber() uint64 {
	block, err := e.client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	log.Println("blockNumber: ", block)
	return block
}

func (e *EthBasicCaller) GetGasPrice() *big.Int {

	price, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	log.Println("GasPrice: ", price)
	return price
}

func (e *EthBasicCaller) GetNonce(address *common.Address) uint64 {

	nonce, err := e.client.PendingNonceAt(context.Background(), *address)
	if err != nil {
		panic(err)
	}
	log.Println(address.Hex(), " nonce is: ", nonce)
	return nonce
}

func (e *EthBasicCaller) GetEstimateGas(from, to *common.Address, input *[]byte) (uint64, error) {

	msg := ethereum.CallMsg{
		From:      *from,
		To:        to,
		GasPrice:  nil,
		GasTipCap: nil,
		GasFeeCap: nil,
		Value:     new(big.Int),
		Data:      *input,
	}

	return e.client.EstimateGas(context.Background(), msg)
}

func (e *EthBasicCaller) GetTransaction(gasPrice *big.Int, nonce, gasLimit uint64, to *common.Address, input *[]byte) (*types.Transaction, error) {

	baseTx := &types.LegacyTx{
		To:       to,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Value:    new(big.Int),
		Data:     *input,
	}

	return types.NewTx(baseTx), nil
}

func (e *EthBasicCaller) SignTransaction(tx *types.Transaction, prv *ecdsa.PrivateKey) *types.Transaction {

	chainId, err := e.client.NetworkID(context.Background())
	if err != nil {
		log.Panicln("ERROR_chainId:", err)
	}
	signedTx, err := SignTx(tx, types.NewEIP155Signer(chainId), prv)
	if err != nil {
		log.Panicln("ERROR_signed:", err)
	}

	return signedTx

}

// SignTx signs the transaction using the given signer and private key.
func SignTx(tx *types.Transaction, s types.Signer, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	log.Println("sign:", hex.EncodeToString(sig))
	return tx.WithSignature(s, sig)
}

func (e *EthBasicCaller) Send(txSigner *types.Transaction) string {
	data, _ := txSigner.MarshalBinary()

	log.Println("send:", hexutil.Encode(data))

	err := e.client.SendTransaction(context.Background(), txSigner)
	if err != nil {
		log.Panicln(err)
	}

	return txSigner.Hash().Hex()
}

func (e *EthBasicCaller) Close() {
	log.Println("Close")
	e.client.Close()
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}
