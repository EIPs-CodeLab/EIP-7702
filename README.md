# EIP-7702 in Go: Core Scaffold + Examples

This repository is a practical Go starter kit for building around **[EIP-7702](https://eips.ethereum.org/EIPS/eip-7702)**.

It focuses on:
- Basic authorization tuple creation/signing/recovery
- Transaction batching payload construction
- Sending ERC-4337 user operations for delegated EOAs

It is designed for EIP workshops / codelabs where you need a clean and hackable foundation.

## What EIP-7702 Enables

EIP-7702 adds a typed transaction (`0x04`) that allows an EOA to authorize delegation to smart-account-like logic. This unlocks:
- Batched execution
- Sponsored gas flows
- Custom permission systems
- Reuse of existing smart account patterns

In short: EOAs can temporarily behave like programmable accounts through delegated code designation.

## Project Goals

- Keep the core primitives explicit and easy to audit
- Keep examples runnable and independent
- Separate reusable package code from demo binaries

## Repository Layout

```text
.
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── release.yml
│       └── vulncheck.yml
├── docs/
│   └── IMPLEMENTATION_NOTES.md
├── examples/
│   ├── basic-authorization/
│   │   └── main.go
│   ├── send-userop/
│   │   └── main.go
│   └── transaction-batching/
│       └── main.go
├── pkg/
│   ├── batching/
│   │   ├── batching.go
│   │   └── batching_test.go
│   ├── eip7702/
│   │   ├── authorization.go
│   │   ├── authorization_test.go
│   │   ├── delegation.go
│   │   ├── delegation_test.go
│   │   ├── doc.go
│   │   ├── setcode_tx.go
│   │   ├── setcode_tx_test.go
│   │   └── types.go
│   └── userop/
│       ├── client.go
│       ├── client_test.go
│       └── types.go
├── scripts/
│   └── build-release.sh
├── Makefile
├── go.mod
└── README.md
```

## Package Overview

### `pkg/eip7702`
Core EIP-7702 helpers:
- Authorization tuple digest: `keccak(0x05 || rlp([chain_id, address, nonce]))`
- Authorization signing and signer recovery
- Low-S signature checks (EIP-2 rule)
- Delegation designation encoding (`0xef0100 || address`)
- Set-code typed transaction payload encoding (`0x04 || rlp([...])`)

### `pkg/batching`
Helpers for batched calls:
- `executeBatch((address,uint256,bytes)[])` calldata encoding
- Generic ABI call encoder for demos

### `pkg/userop`
Small JSON-RPC bundler client for ERC-4337:
- UserOperation struct
- Request builder for `eth_sendUserOperation`
- HTTP client for submission

## Quick Start

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Run tests

```bash
go test ./...
```

### 3. Run examples

```bash
go run ./examples/basic-authorization
go run ./examples/transaction-batching
go run ./examples/send-userop
```

### 4. Run full local CI

```bash
make ci
```

## CI/CD

This project includes three GitHub Actions workflows:

- **CI** (`.github/workflows/ci.yml`)
  - Triggered on pull requests and pushes to `main` / `master`
  - Runs format checks, module tidy check, vet, race tests, build, and example runs
- **Release CD** (`.github/workflows/release.yml`)
  - Triggered on tags matching `v*` (for example `v0.1.0`)
  - Builds cross-platform example binaries and publishes release artifacts
- **Vulnerability Scan** (`.github/workflows/vulncheck.yml`)
  - Weekly scheduled `govulncheck` + manual run support

### Release process

1. Push a semantic tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

2. The release workflow will:
  - Build artifacts for `linux`, `darwin`, and `windows`
  - Generate `dist/checksums.txt`
  - Publish assets to a GitHub Release

## Example Walkthroughs

### A) Basic Authorization

```bash
go run ./examples/basic-authorization
```

What it demonstrates:
- Creates an authorization tuple
- Signs with secp256k1 key
- Recovers authority address from signature
- Builds delegation code marker

### B) Transaction Batching

```bash
go run ./examples/transaction-batching
```

What it demonstrates:
- Encodes multiple ERC-20 transfers into one batch call
- Adds authorization list
- Builds EIP-7702 typed tx bytes (`0x04...`)

Note: The example intentionally leaves outer transaction signing as placeholder so the flow is easy to inspect.

### C) Sending a UserOperation

```bash
go run ./examples/send-userop
```

Default behavior is dry-run (prints JSON-RPC request).

To broadcast:

```bash
export BUNDLER_RPC_URL="https://your-bundler-rpc"
export ENTRYPOINT="0x0000000071727De22E5E9d8BAf0edAc6f37da032" # optional

go run ./examples/send-userop
```

## Implementation Notes

- The code models EIP-7702 wire formats and authorization logic for application-layer tooling.
- This repository does **not** ship a full execution client or consensus-level integration.
- Validate chain-specific constants and RPC behavior in your target environment.

Read more in [`docs/IMPLEMENTATION_NOTES.md`](docs/IMPLEMENTATION_NOTES.md).

## Security Checklist Before Production

- Use hardware-backed key management / signer services
- Enforce nonce management and replay controls per chain
- Verify delegate contract semantics and upgrade policies
- Verify paymaster and bundler trust boundaries
- Add simulation + revert reason inspection before submission

## References

- EIP-7702: [https://eips.ethereum.org/EIPS/eip-7702](https://eips.ethereum.org/EIPS/eip-7702)
- Ethereum Magicians thread: [https://ethereum-magicians.org/t/eip-set-eoa-account-code-for-one-transaction/19923](https://ethereum-magicians.org/t/eip-set-eoa-account-code-for-one-transaction/19923)
- Proxy pattern gist: [https://gist.github.com/lightclient/7742e84fde4962f32928c6177eda7523](https://gist.github.com/lightclient/7742e84fde4962f32928c6177eda7523)
- Best practices draft: [https://hackmd.io/@rimeissner/eip7702-best-practices](https://hackmd.io/@rimeissner/eip7702-best-practices)
