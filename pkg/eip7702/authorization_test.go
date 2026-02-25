package eip7702_test

import (
	"math/big"
	"testing"

	"github.com/eipcodelab/eip7702-go/pkg/eip7702"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSignAndRecoverAuthorization(t *testing.T) {
	key, err := crypto.HexToECDSA("4f3edf983ac63f7f8b7d0c4f76f2a5a70fadb53fcbf65f45d6fd5d77f07683ab")
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	delegate := common.HexToAddress("0x000000000000000000000000000000000000c0de")
	auth, err := eip7702.SignAuthorization(key, big.NewInt(1), delegate, 0)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	recovered, err := eip7702.VerifyAuthorization(auth, big.NewInt(1))
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	want := crypto.PubkeyToAddress(key.PublicKey)
	if recovered != want {
		t.Fatalf("unexpected signer: got %s want %s", recovered.Hex(), want.Hex())
	}
}

func TestVerifyAuthorizationRejectsWrongChain(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	auth, err := eip7702.SignAuthorization(key, big.NewInt(10), common.Address{}, 1)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := eip7702.VerifyAuthorization(auth, big.NewInt(1)); err == nil {
		t.Fatal("expected chain mismatch error")
	}
}

func TestVerifyAuthorizationAcceptsChainZero(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	auth, err := eip7702.SignAuthorization(key, big.NewInt(0), common.Address{}, 2)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := eip7702.VerifyAuthorization(auth, big.NewInt(1)); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestVerifyAuthorizationRejectsHighS(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	auth, err := eip7702.SignAuthorization(key, big.NewInt(1), common.Address{}, 0)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	auth.S = new(big.Int).Sub(crypto.S256().Params().N, auth.S)
	if auth.S.Cmp(new(big.Int).Rsh(crypto.S256().Params().N, 1)) <= 0 {
		t.Fatal("expected high-S for test setup")
	}
	if _, err := eip7702.VerifyAuthorization(auth, big.NewInt(1)); err == nil {
		t.Fatal("expected low-S validation error")
	}
}
