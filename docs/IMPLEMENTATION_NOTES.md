# EIP-7702 Implementation Notes (Go)

This document describes how this scaffold maps the EIP into code.

## 1. Authorization Tuple

The EIP authorization item:

```text
[chain_id, address, nonce, y_parity, r, s]
```

In code:
- `pkg/eip7702/types.go` defines `Authorization`
- `pkg/eip7702/authorization.go` computes the digest:
  - `keccak(0x05 || rlp([chain_id, address, nonce]))`
- Signing returns `y_parity`, `r`, `s` in RLP-ready form
- Verification includes:
  - chain-id compatibility (`0` or current chain)
  - low-S check
  - signer recovery

## 2. Delegation Designation

EIP-7702 writes delegated code as:

```text
0xef0100 || address
```

In code:
- `pkg/eip7702/delegation.go`
  - `DelegationCode(delegate)`
  - `ParseDelegationCode(code)`
  - `IsClearCodeAuthorization(delegate)` for `0x0` address case

## 3. Typed Transaction Encoding

Set-code tx payload is encoded in strict field order using RLP.

In code:
- `pkg/eip7702/types.go` defines `SetCodeTx`
- `pkg/eip7702/setcode_tx.go`
  - `EncodePayload()`
  - `EncodeTypedTransaction()` => `0x04 || payload`

## 4. Batching Strategy

Batching is demonstrated using an `executeBatch((address,uint256,bytes)[])` ABI pattern.

In code:
- `pkg/batching/batching.go`
  - `EncodeExecuteBatch(calls)`
  - `EncodeFunctionCall(...)`

## 5. UserOperation Submission

For EIP-4337 compatibility examples:
- `pkg/userop/types.go` defines `UserOperation`
- `pkg/userop/client.go` builds and sends `eth_sendUserOperation`

The example program (`examples/send-userop/main.go`) prints payload by default and submits only when `BUNDLER_RPC_URL` is set.

## 6. Scope Boundaries

This repository intentionally avoids client-internal state transition logic and consensus rules. It focuses on:
- App-layer encoding/signing
- Interop with wallet/bundler flows
- Testable primitives for codelab usage

For production, pair this with:
- chain-specific transaction signing support for type `0x04`
- simulation infrastructure
- robust key management
