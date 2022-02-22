package nat

import (
	"encoding/hex"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestTransferByteCode(t *testing.T) {
	api, _ := BuildApi()

	to := common.HexToAddress("0xB3cedeebaea9F8DeA537cc8e6136B6338104FE0f")

	amount := big.NewInt(1)
	byteCode, _ := TransferByteCode(api, &to, amount)

	result := hex.EncodeToString(*byteCode)
	log.Println(result)

	if result != "a9059cbb000000000000000000000000d69399c5f03762315946dd44badcdec2451d66330000000000000000000000000000000000000000000000000000000000000064" {
		t.Error("bytecode incorrect")
	}

}
