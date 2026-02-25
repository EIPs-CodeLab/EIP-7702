package eip7702

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

// EncodePayload returns the RLP payload of a type-0x04 transaction.
func (tx *SetCodeTx) EncodePayload() ([]byte, error) {
	if err := tx.ValidateBasic(); err != nil {
		return nil, err
	}
	enc, err := rlp.EncodeToBytes([]any{
		tx.ChainID,
		tx.Nonce,
		tx.MaxPriorityFeePerGas,
		tx.MaxFeePerGas,
		tx.GasLimit,
		tx.Destination,
		tx.Value,
		tx.Data,
		tx.AccessList,
		tx.AuthorizationList,
		tx.SignatureYParity,
		tx.SignatureR,
		tx.SignatureS,
	})
	if err != nil {
		return nil, fmt.Errorf("encode set-code payload: %w", err)
	}
	return enc, nil
}

// EncodeTypedTransaction prefixes 0x04 to the payload as required by EIP-2718.
func (tx *SetCodeTx) EncodeTypedTransaction() ([]byte, error) {
	payload, err := tx.EncodePayload()
	if err != nil {
		return nil, err
	}
	out := make([]byte, 1+len(payload))
	out[0] = SetCodeTxType
	copy(out[1:], payload)
	return out, nil
}
