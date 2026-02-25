package eip7702_test

import (
	"math/big"
	"testing"

	"github.com/eipcodelab/eip7702-go/pkg/eip7702"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSetCodeTxEncoding(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	auth, err := eip7702.SignAuthorization(key, big.NewInt(1), common.HexToAddress("0x2000000000000000000000000000000000000002"), 0)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	tx := &eip7702.SetCodeTx{
		ChainID:              big.NewInt(1),
		Nonce:                0,
		MaxPriorityFeePerGas: big.NewInt(2_000_000_000),
		MaxFeePerGas:         big.NewInt(40_000_000_000),
		GasLimit:             250_000,
		Destination:          common.HexToAddress("0x3000000000000000000000000000000000000003"),
		Value:                big.NewInt(0),
		Data:                 []byte{0xde, 0xad, 0xbe, 0xef},
		AuthorizationList:    []eip7702.Authorization{auth},
		SignatureYParity:     0,
		SignatureR:           big.NewInt(1),
		SignatureS:           big.NewInt(1),
	}

	raw, err := tx.EncodeTypedTransaction()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(raw) < 2 {
		t.Fatal("typed tx encoding is too short")
	}
	if raw[0] != eip7702.SetCodeTxType {
		t.Fatalf("unexpected tx type: got 0x%x", raw[0])
	}
}
