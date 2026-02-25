package eip7702

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

var (
	delegationPrefix = []byte{0xef, 0x01, 0x00}
	// EmptyCodeHash is keccak256("") used when clearing delegated code.
	EmptyCodeHash = common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
)

// DelegationCode returns 0xef0100 || delegate. Zero address means "clear code".
func DelegationCode(delegate common.Address) []byte {
	if IsClearCodeAuthorization(delegate) {
		return nil
	}
	out := make([]byte, len(delegationPrefix)+common.AddressLength)
	copy(out, delegationPrefix)
	copy(out[len(delegationPrefix):], delegate.Bytes())
	return out
}

// ParseDelegationCode parses code designation and returns the delegated address.
func ParseDelegationCode(code []byte) (common.Address, bool) {
	if len(code) != len(delegationPrefix)+common.AddressLength {
		return common.Address{}, false
	}
	if !bytes.Equal(code[:len(delegationPrefix)], delegationPrefix) {
		return common.Address{}, false
	}
	var target common.Address
	copy(target[:], code[len(delegationPrefix):])
	return target, true
}

// IsClearCodeAuthorization is true when tuple address is 0x0.
func IsClearCodeAuthorization(delegate common.Address) bool {
	return delegate == (common.Address{})
}
