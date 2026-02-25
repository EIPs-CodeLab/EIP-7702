package batching_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/eipcodelab/eip7702-go/pkg/batching"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestEncodeExecuteBatch(t *testing.T) {
	calls := []batching.Call{
		{Target: common.HexToAddress("0x0000000000000000000000000000000000000001"), Value: big.NewInt(0), Data: []byte{0x01}},
		{Target: common.HexToAddress("0x0000000000000000000000000000000000000002"), Value: big.NewInt(1), Data: []byte{0x02}},
	}
	calldata, err := batching.EncodeExecuteBatch(calls)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(calldata) < 4 {
		t.Fatal("calldata is too short")
	}
	methodID := crypto.Keccak256([]byte("executeBatch((address,uint256,bytes)[])"))[:4]
	if !bytes.Equal(methodID, calldata[:4]) {
		t.Fatalf("unexpected selector: got %x want %x", calldata[:4], methodID)
	}
}

func TestEncodeExecuteBatchRejectsEmptyCalls(t *testing.T) {
	if _, err := batching.EncodeExecuteBatch(nil); err == nil {
		t.Fatal("expected error for empty call list")
	}
}
