package batching

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const executeBatchABIJSON = `[
  {
    "type": "function",
    "name": "executeBatch",
    "stateMutability": "payable",
    "inputs": [
      {
        "name": "calls",
        "type": "tuple[]",
        "components": [
          {"name": "target", "type": "address"},
          {"name": "value", "type": "uint256"},
          {"name": "data", "type": "bytes"}
        ]
      }
    ],
    "outputs": []
  }
]`

type executeCall struct {
	Target common.Address `abi:"target"`
	Value  *big.Int       `abi:"value"`
	Data   []byte         `abi:"data"`
}

// Call represents one low-level call executed by a batch delegate contract.
type Call struct {
	Target common.Address
	Value  *big.Int
	Data   []byte
}

var (
	onceBatchABI sync.Once
	batchABI     abi.ABI
	batchABIErr  error
)

func getExecuteBatchABI() (abi.ABI, error) {
	onceBatchABI.Do(func() {
		batchABI, batchABIErr = abi.JSON(strings.NewReader(executeBatchABIJSON))
	})
	return batchABI, batchABIErr
}

// EncodeExecuteBatch encodes calldata for executeBatch((address,uint256,bytes)[] calls).
func EncodeExecuteBatch(calls []Call) ([]byte, error) {
	if len(calls) == 0 {
		return nil, errors.New("calls must not be empty")
	}
	execCalls := make([]executeCall, len(calls))
	for i, c := range calls {
		value := c.Value
		if value == nil {
			value = big.NewInt(0)
		}
		execCalls[i] = executeCall{
			Target: c.Target,
			Value:  value,
			Data:   c.Data,
		}
	}

	parsedABI, err := getExecuteBatchABI()
	if err != nil {
		return nil, fmt.Errorf("parse executeBatch ABI: %w", err)
	}
	data, err := parsedABI.Pack("executeBatch", execCalls)
	if err != nil {
		return nil, fmt.Errorf("pack executeBatch calldata: %w", err)
	}
	return data, nil
}

// EncodeFunctionCall is a small helper for example programs.
func EncodeFunctionCall(abiJSON string, method string, args ...any) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("parse ABI: %w", err)
	}
	out, err := parsedABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("pack %s: %w", method, err)
	}
	return out, nil
}
