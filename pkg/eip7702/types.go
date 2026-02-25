package eip7702

import (
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// SetCodeTxType is the typed transaction ID introduced by EIP-7702.
	SetCodeTxType byte = 0x04
	// AuthorizationMagic is prepended to the authorization tuple before hashing.
	AuthorizationMagic byte = 0x05

	// PER_AUTH_BASE_COST is the base gas charged for each authorization tuple.
	PER_AUTH_BASE_COST uint64 = 12_500
	// PER_EMPTY_ACCOUNT_COST is the warm account cost used for refund accounting.
	PER_EMPTY_ACCOUNT_COST uint64 = 25_000
)

var (
	ErrNilChainID        = errors.New("chain id is required")
	ErrInvalidChainID    = errors.New("chain id must be >= 0")
	ErrMaxNonce          = errors.New("nonce must be < 2^64 - 1")
	ErrEmptyAuthList     = errors.New("authorization_list must not be empty")
	ErrInvalidYParity    = errors.New("y parity must be 0 or 1")
	ErrNilSignatureValue = errors.New("signature r/s values are required")
	ErrInvalidSignature  = errors.New("signature r/s must be positive 256-bit values")
)

// Authorization is one item in authorization_list.
// RLP field order follows the EIP tuple definition.
type Authorization struct {
	ChainID *big.Int       `json:"chainId"`
	Address common.Address `json:"address"`
	Nonce   uint64         `json:"nonce"`
	YParity uint8          `json:"yParity"`
	R       *big.Int       `json:"r"`
	S       *big.Int       `json:"s"`
}

// ValidateBasic enforces tuple-level checks from the EIP that can be done offline.
func (a Authorization) ValidateBasic() error {
	if a.ChainID == nil {
		return ErrNilChainID
	}
	if a.ChainID.Sign() < 0 {
		return ErrInvalidChainID
	}
	if a.Nonce == math.MaxUint64 {
		return ErrMaxNonce
	}
	if a.YParity > 1 {
		return ErrInvalidYParity
	}
	if a.R == nil || a.S == nil {
		return ErrNilSignatureValue
	}
	if a.R.Sign() <= 0 || a.S.Sign() <= 0 {
		return ErrInvalidSignature
	}
	if a.R.BitLen() > 256 || a.S.BitLen() > 256 {
		return ErrInvalidSignature
	}
	return nil
}

// SetCodeTx models the EIP-7702 typed transaction payload.
type SetCodeTx struct {
	ChainID              *big.Int         `json:"chainId"`
	Nonce                uint64           `json:"nonce"`
	MaxPriorityFeePerGas *big.Int         `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *big.Int         `json:"maxFeePerGas"`
	GasLimit             uint64           `json:"gasLimit"`
	Destination          common.Address   `json:"destination"`
	Value                *big.Int         `json:"value"`
	Data                 []byte           `json:"data"`
	AccessList           types.AccessList `json:"accessList"`
	AuthorizationList    []Authorization  `json:"authorizationList"`
	SignatureYParity     uint8            `json:"signatureYParity"`
	SignatureR           *big.Int         `json:"signatureR"`
	SignatureS           *big.Int         `json:"signatureS"`
}

// ValidateBasic validates required set-code fields before encoding/signing.
func (tx *SetCodeTx) ValidateBasic() error {
	if tx.ChainID == nil {
		return ErrNilChainID
	}
	if tx.ChainID.Sign() < 0 {
		return ErrInvalidChainID
	}
	if tx.Nonce == math.MaxUint64 {
		return ErrMaxNonce
	}
	if len(tx.AuthorizationList) == 0 {
		return ErrEmptyAuthList
	}
	for _, auth := range tx.AuthorizationList {
		if err := auth.ValidateBasic(); err != nil {
			return err
		}
	}
	if tx.MaxPriorityFeePerGas == nil || tx.MaxFeePerGas == nil || tx.Value == nil {
		return errors.New("maxPriorityFeePerGas, maxFeePerGas and value are required")
	}
	if tx.SignatureYParity > 1 {
		return ErrInvalidYParity
	}
	if tx.SignatureR == nil || tx.SignatureS == nil {
		return ErrNilSignatureValue
	}
	if tx.SignatureR.Sign() <= 0 || tx.SignatureS.Sign() <= 0 {
		return ErrInvalidSignature
	}
	if tx.SignatureR.BitLen() > 256 || tx.SignatureS.BitLen() > 256 {
		return ErrInvalidSignature
	}
	return nil
}

// AuthorizationRefundDelta returns the refund increment used by EIP-7702 when authority exists.
func AuthorizationRefundDelta() uint64 {
	return PER_EMPTY_ACCOUNT_COST - PER_AUTH_BASE_COST
}
