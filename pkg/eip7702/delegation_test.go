package eip7702_test

import (
	"testing"

	"github.com/eipcodelab/eip7702-go/pkg/eip7702"
	"github.com/ethereum/go-ethereum/common"
)

func TestDelegationCodeRoundtrip(t *testing.T) {
	delegate := common.HexToAddress("0x1000000000000000000000000000000000000001")
	code := eip7702.DelegationCode(delegate)
	parsed, ok := eip7702.ParseDelegationCode(code)
	if !ok {
		t.Fatal("expected valid delegation code")
	}
	if parsed != delegate {
		t.Fatalf("unexpected delegate: got %s want %s", parsed.Hex(), delegate.Hex())
	}
}

func TestDelegationCodeClearFlow(t *testing.T) {
	if out := eip7702.DelegationCode(common.Address{}); out != nil {
		t.Fatal("zero delegate must return nil code designation")
	}
	if !eip7702.IsClearCodeAuthorization(common.Address{}) {
		t.Fatal("zero address should be treated as clear-code")
	}
}
