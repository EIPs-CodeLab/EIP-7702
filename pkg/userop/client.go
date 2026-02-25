package userop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type rpcError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

// BundlerClient is a tiny JSON-RPC client for ERC-4337 calls.
type BundlerClient struct {
	endpoint   string
	httpClient *http.Client
}

// NewBundlerClient creates a client for eth_sendUserOperation requests.
func NewBundlerClient(endpoint string) *BundlerClient {
	return &BundlerClient{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// NewBundlerClientWithHTTPClient is useful for tests and custom transport wiring.
func NewBundlerClientWithHTTPClient(endpoint string, httpClient *http.Client) *BundlerClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 20 * time.Second}
	}
	return &BundlerClient{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

// BuildSendUserOperationRequest returns JSON payload for eth_sendUserOperation.
func BuildSendUserOperationRequest(op UserOperation, entryPoint common.Address) ([]byte, error) {
	if err := op.ValidateBasic(); err != nil {
		return nil, err
	}
	body := rpcRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_sendUserOperation",
		Params:  []any{op, entryPoint},
	}
	enc, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal rpc request: %w", err)
	}
	return enc, nil
}

// SendUserOperation sends one user operation and returns the userOp hash.
func (c *BundlerClient) SendUserOperation(ctx context.Context, op UserOperation, entryPoint common.Address) (common.Hash, error) {
	payload, err := BuildSendUserOperationRequest(op, entryPoint)
	if err != nil {
		return common.Hash{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return common.Hash{}, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return common.Hash{}, fmt.Errorf("post request: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return common.Hash{}, fmt.Errorf("decode rpc response: %w", err)
	}
	if rpcResp.Error != nil {
		return common.Hash{}, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	var hash common.Hash
	if err := json.Unmarshal(rpcResp.Result, &hash); err != nil {
		return common.Hash{}, fmt.Errorf("decode result hash: %w", err)
	}
	return hash, nil
}
