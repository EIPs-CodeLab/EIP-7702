package userop_test

import (
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"

	"github.com/eipcodelab/eip7702-go/pkg/userop"
	"github.com/ethereum/go-ethereum/common"
)

func makeUserOp() userop.UserOperation {
	return userop.UserOperation{
		Sender:               common.HexToAddress("0x000000000000000000000000000000000000dead"),
		Nonce:                userop.HexBig(big.NewInt(1)),
		CallData:             []byte{0x01, 0x02},
		CallGasLimit:         userop.HexUint64(100_000),
		VerificationGasLimit: userop.HexUint64(250_000),
		PreVerificationGas:   userop.HexUint64(50_000),
		MaxFeePerGas:         userop.HexBig(big.NewInt(30_000_000_000)),
		MaxPriorityFeePerGas: userop.HexBig(big.NewInt(2_000_000_000)),
		Signature:            []byte{0xaa, 0xbb},
	}
}

func TestBuildSendUserOperationRequest(t *testing.T) {
	op := makeUserOp()
	payload, err := userop.BuildSendUserOperationRequest(op, common.HexToAddress("0x0000000000000000000000000000000000000001"))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(payload, &body); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if body["method"] != "eth_sendUserOperation" {
		t.Fatalf("unexpected method: %v", body["method"])
	}
}

func TestSendUserOperation(t *testing.T) {
	client := userop.NewBundlerClientWithHTTPClient(
		"https://bundler.example",
		&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != "https://bundler.example" {
					t.Fatalf("unexpected endpoint: %s", req.URL.String())
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				if !strings.Contains(string(body), "eth_sendUserOperation") {
					t.Fatalf("unexpected rpc body: %s", string(body))
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body: io.NopCloser(strings.NewReader(
						`{"jsonrpc":"2.0","id":1,"result":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`,
					)),
				}, nil
			}),
		},
	)
	hash, err := client.SendUserOperation(context.Background(), makeUserOp(), common.HexToAddress("0x0000000000000000000000000000000000000001"))
	if err != nil {
		t.Fatalf("send user op: %v", err)
	}
	if hash.Hex() != "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("unexpected hash: %s", hash.Hex())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
