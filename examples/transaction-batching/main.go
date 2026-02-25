package main

import (
	"fmt"
	"math/big"

	"github.com/eipcodelab/eip7702-go/pkg/batching"
	"github.com/eipcodelab/eip7702-go/pkg/eip7702"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const erc20ABI = `[
  {
    "type": "function",
    "name": "transfer",
    "stateMutability": "nonpayable",
    "inputs": [
      {"name": "to", "type": "address"},
      {"name": "value", "type": "uint256"}
    ],
    "outputs": [{"name": "", "type": "bool"}]
  }
]`

func main() {
	key, err := crypto.HexToECDSA("4f3edf983ac63f7f8b7d0c4f76f2a5a70fadb53fcbf65f45d6fd5d77f07683ab")
	if err != nil {
		panic(err)
	}
	authority := crypto.PubkeyToAddress(key.PublicKey)

	chainID := big.NewInt(1)
	delegate := common.HexToAddress("0x1111111111111111111111111111111111111111")

	transferA, err := batching.EncodeFunctionCall(
		erc20ABI,
		"transfer",
		common.HexToAddress("0x2000000000000000000000000000000000000002"),
		big.NewInt(1_000_000_000_000_000_000), // 1 token with 18 decimals
	)
	if err != nil {
		panic(err)
	}
	transferB, err := batching.EncodeFunctionCall(
		erc20ABI,
		"transfer",
		common.HexToAddress("0x3000000000000000000000000000000000000003"),
		big.NewInt(250_000_000_000_000_000),
	)
	if err != nil {
		panic(err)
	}

	batchCalldata, err := batching.EncodeExecuteBatch([]batching.Call{
		{
			Target: common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // example ERC20
			Value:  big.NewInt(0),
			Data:   transferA,
		},
		{
			Target: common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
			Value:  big.NewInt(0),
			Data:   transferB,
		},
	})
	if err != nil {
		panic(err)
	}

	auth, err := eip7702.SignAuthorization(key, chainID, delegate, 0)
	if err != nil {
		panic(err)
	}

	setCodeTx := &eip7702.SetCodeTx{
		ChainID:              chainID,
		Nonce:                0,
		MaxPriorityFeePerGas: big.NewInt(2_000_000_000),
		MaxFeePerGas:         big.NewInt(40_000_000_000),
		GasLimit:             300_000,
		Destination:          authority,
		Value:                big.NewInt(0),
		Data:                 batchCalldata,
		AuthorizationList:    []eip7702.Authorization{auth},
		// Outer tx signature is placeholder in this example.
		SignatureYParity: 0,
		SignatureR:       big.NewInt(1),
		SignatureS:       big.NewInt(1),
	}

	raw, err := setCodeTx.EncodeTypedTransaction()
	if err != nil {
		panic(err)
	}

	fmt.Println("== EIP-7702 Transaction Batching ==")
	fmt.Printf("Authority:           %s\n", authority.Hex())
	fmt.Printf("Delegation target:   %s\n", delegate.Hex())
	fmt.Printf("Batch calldata:      0x%x\n", batchCalldata)
	fmt.Printf("Typed tx (0x04...):  0x%x\n", raw)
	fmt.Println("Note: sign the outer set-code transaction before broadcasting.")
}
