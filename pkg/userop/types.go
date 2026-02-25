package userop

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// UserOperation follows the ERC-4337 v0.6 JSON shape used by eth_sendUserOperation.
type UserOperation struct {
	Sender               common.Address `json:"sender"`
	Nonce                *hexutil.Big   `json:"nonce"`
	InitCode             hexutil.Bytes  `json:"initCode"`
	CallData             hexutil.Bytes  `json:"callData"`
	CallGasLimit         hexutil.Uint64 `json:"callGasLimit"`
	VerificationGasLimit hexutil.Uint64 `json:"verificationGasLimit"`
	PreVerificationGas   hexutil.Uint64 `json:"preVerificationGas"`
	MaxFeePerGas         *hexutil.Big   `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big   `json:"maxPriorityFeePerGas"`
	PaymasterAndData     hexutil.Bytes  `json:"paymasterAndData"`
	Signature            hexutil.Bytes  `json:"signature"`
}

// ValidateBasic applies local sanity checks before sending to a bundler.
func (u UserOperation) ValidateBasic() error {
	if u.Nonce == nil {
		return errors.New("nonce is required")
	}
	if u.MaxFeePerGas == nil || u.MaxPriorityFeePerGas == nil {
		return errors.New("max fee fields are required")
	}
	if len(u.CallData) == 0 {
		return errors.New("callData must not be empty")
	}
	if len(u.Signature) == 0 {
		return errors.New("signature must not be empty")
	}
	return nil
}

// HexBig converts big.Int values into JSON-RPC hex quantities.
func HexBig(v *big.Int) *hexutil.Big {
	if v == nil {
		return (*hexutil.Big)(big.NewInt(0))
	}
	return (*hexutil.Big)(new(big.Int).Set(v))
}

// HexUint64 converts uint64 values into JSON-RPC hex quantities.
func HexUint64(v uint64) hexutil.Uint64 {
	return hexutil.Uint64(v)
}
