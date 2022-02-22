package lib

import (
	"log"
	"testing"
)

func TestPublicToAddress(t *testing.T) {
	pub := "1386edc2cb0173ea7ac4fef74185970521bb225af9316917ad5a80ee65ff8b870f09b4e57665c30539214cf447f23b732e4e32fa4bcf0f710ef81a5c2e4c19c7"

	address := PublicToAddress(pub)

	log.Println("final address:", address.Hex())

}

func TestPrivateToAddress(t *testing.T) {
	pk := "6def68041b8e7de549a549f3daa0573628b90f74871945fbe72940f7e7745b1a"

	address := PrivateToAddress(pk)

	log.Println("final address:", address.Hex())

}
