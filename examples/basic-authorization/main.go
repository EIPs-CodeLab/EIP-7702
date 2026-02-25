package main

import (
	"fmt"
	"math/big"

	"github.com/eipcodelab/eip7702-go/pkg/eip7702"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Demo key used in local dev tooling only.
	key, err := crypto.HexToECDSA("4f3edf983ac63f7f8b7d0c4f76f2a5a70fadb53fcbf65f45d6fd5d77f07683ab")
	if err != nil {
		panic(err)
	}

	chainID := big.NewInt(1)
	delegate := common.HexToAddress("0x1111111111111111111111111111111111111111")
	nonce := uint64(0)

	auth, err := eip7702.SignAuthorization(key, chainID, delegate, nonce)
	if err != nil {
		panic(err)
	}
	authority, err := eip7702.VerifyAuthorization(auth, chainID)
	if err != nil {
		panic(err)
	}
	digest, err := eip7702.AuthorizationDigest(chainID, delegate, nonce)
	if err != nil {
		panic(err)
	}

	fmt.Println("== Basic EIP-7702 Authorization ==")
	fmt.Printf("Authority (signer):   %s\n", authority.Hex())
	fmt.Printf("Delegate target:      %s\n", delegate.Hex())
	fmt.Printf("Authorization digest: 0x%x\n", digest)
	fmt.Printf("Tuple: [chain_id=%s, address=%s, nonce=%d, y_parity=%d, r=0x%x, s=0x%x]\n",
		auth.ChainID.String(), auth.Address.Hex(), auth.Nonce, auth.YParity, auth.R, auth.S)

	delegationCode := eip7702.DelegationCode(delegate)
	fmt.Printf("Code designation:     0x%x\n", delegationCode)
}
