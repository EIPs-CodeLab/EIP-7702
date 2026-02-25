package eip7702

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var secp256k1HalfN = new(big.Int).Rsh(new(big.Int).Set(crypto.S256().Params().N), 1)

// AuthorizationDigest computes keccak(0x05 || rlp([chain_id, address, nonce])).
func AuthorizationDigest(chainID *big.Int, target common.Address, nonce uint64) ([]byte, error) {
	if chainID == nil {
		return nil, ErrNilChainID
	}
	encoded, err := rlp.EncodeToBytes([]any{chainID, target, nonce})
	if err != nil {
		return nil, fmt.Errorf("encode tuple: %w", err)
	}
	payload := append([]byte{AuthorizationMagic}, encoded...)
	return crypto.Keccak256(payload), nil
}

// SignAuthorization signs an authorization tuple and returns an RLP-ready struct.
func SignAuthorization(privateKey *ecdsa.PrivateKey, chainID *big.Int, delegate common.Address, nonce uint64) (Authorization, error) {
	if privateKey == nil {
		return Authorization{}, errors.New("private key is required")
	}
	digest, err := AuthorizationDigest(chainID, delegate, nonce)
	if err != nil {
		return Authorization{}, err
	}
	sig, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return Authorization{}, fmt.Errorf("sign tuple: %w", err)
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	if s.Cmp(secp256k1HalfN) > 0 {
		return Authorization{}, errors.New("signature S is not low-S")
	}

	auth := Authorization{
		ChainID: chainID,
		Address: delegate,
		Nonce:   nonce,
		YParity: sig[64],
		R:       r,
		S:       s,
	}
	if err := auth.ValidateBasic(); err != nil {
		return Authorization{}, err
	}
	return auth, nil
}

// RecoverAuthority recovers the signer address from one authorization tuple.
func RecoverAuthority(auth Authorization) (common.Address, error) {
	if err := auth.ValidateBasic(); err != nil {
		return common.Address{}, err
	}
	digest, err := AuthorizationDigest(auth.ChainID, auth.Address, auth.Nonce)
	if err != nil {
		return common.Address{}, err
	}
	sig := make([]byte, 65)
	rBytes := auth.R.Bytes()
	sBytes := auth.S.Bytes()
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)
	sig[64] = auth.YParity

	pub, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return common.Address{}, fmt.Errorf("recover signer: %w", err)
	}
	return crypto.PubkeyToAddress(*pub), nil
}

// VerifyAuthorization applies chain-id and low-S checks and returns recovered authority.
func VerifyAuthorization(auth Authorization, currentChainID *big.Int) (common.Address, error) {
	if currentChainID == nil {
		return common.Address{}, ErrNilChainID
	}
	if err := auth.ValidateBasic(); err != nil {
		return common.Address{}, err
	}
	if auth.ChainID.Sign() != 0 && auth.ChainID.Cmp(currentChainID) != 0 {
		return common.Address{}, errors.New("authorization chain id does not match current chain")
	}
	if auth.S.Cmp(secp256k1HalfN) > 0 {
		return common.Address{}, errors.New("authorization signature violates low-S rule")
	}
	return RecoverAuthority(auth)
}
