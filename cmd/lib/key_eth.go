package lib

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func PublicToAddress(pub string) *common.Address {
	pub = "04" + pub

	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		log.Panicln("Error pu hexa:", err)
		panic(err)
	}

	pubkey1, err := crypto.UnmarshalPubkey(pubBytes)
	if err != nil {
		log.Panicln("Error descompress:", err)
		panic(err)
	}

	address := crypto.PubkeyToAddress(*pubkey1)

	return &address
}

func PrivateToAddress(pk string) *common.Address {

	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &fromAddress
}
