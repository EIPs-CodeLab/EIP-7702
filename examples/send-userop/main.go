package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/eipcodelab/eip7702-go/pkg/batching"
	"github.com/eipcodelab/eip7702-go/pkg/userop"
	"github.com/ethereum/go-ethereum/common"
)

const walletDelegateABI = `[
  {
    "type": "function",
    "name": "execute",
    "stateMutability": "nonpayable",
    "inputs": [
      {"name": "target", "type": "address"},
      {"name": "value", "type": "uint256"},
      {"name": "data", "type": "bytes"}
    ],
    "outputs": []
  }
]`

func main() {
	sender := common.HexToAddress("0x000000000000000000000000000000000000dEaD")
	entryPoint := common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032") // v0.6 example
	if env := os.Getenv("ENTRYPOINT"); env != "" {
		entryPoint = common.HexToAddress(env)
	}

	callData, err := batching.EncodeFunctionCall(
		walletDelegateABI,
		"execute",
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		big.NewInt(0),
		[]byte{},
	)
	if err != nil {
		panic(err)
	}

	op := userop.UserOperation{
		Sender:               sender,
		Nonce:                userop.HexBig(big.NewInt(0)),
		InitCode:             []byte{},
		CallData:             callData,
		CallGasLimit:         userop.HexUint64(200_000),
		VerificationGasLimit: userop.HexUint64(300_000),
		PreVerificationGas:   userop.HexUint64(80_000),
		MaxFeePerGas:         userop.HexBig(big.NewInt(35_000_000_000)),
		MaxPriorityFeePerGas: userop.HexBig(big.NewInt(2_000_000_000)),
		PaymasterAndData:     []byte{},
		Signature:            []byte{0xde, 0xad, 0xbe, 0xef}, // replace with actual signature
	}

	payload, err := userop.BuildSendUserOperationRequest(op, entryPoint)
	if err != nil {
		panic(err)
	}

	fmt.Println("== ERC-4337 UserOperation over EIP-7702 account ==")
	fmt.Println(string(payload))

	endpoint := os.Getenv("BUNDLER_RPC_URL")
	if endpoint == "" {
		fmt.Println("Dry-run only. Set BUNDLER_RPC_URL to broadcast.")
		return
	}

	client := userop.NewBundlerClient(endpoint)
	hash, err := client.SendUserOperation(context.Background(), op, entryPoint)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Submitted userOp hash: %s\n", hash.Hex())
}
