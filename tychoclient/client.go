package tychoclient

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/tonkeeper/tongo/tychoclient/proto"
)

const (
	// DefaultEndpoint is the testnet endpoint provided by the team
	DefaultEndpoint = "tonapi-testnet.tychoprotocol.com:443"

	// DefaultTimeout for gRPC calls
	DefaultTimeout = 30 * time.Second
)

var BlockNotFoundErr = errors.New("block not found")

// Client provides access to Tycho blockchain data via gRPC
type Client struct {
	conn   *grpc.ClientConn
	client proto.TychoIndexerClient
}

// NewClient creates a new Tycho client with default settings
func NewClient() (*Client, error) {
	return NewClientWithEndpoint(DefaultEndpoint)
}

// NewClientWithEndpoint creates a new Tycho client with custom endpoint
func NewClientWithEndpoint(endpoint string) (*Client, error) {
	// Use TLS for secure connection
	creds := credentials.NewTLS(&tls.Config{})

	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", endpoint, err)
	}

	return &Client{
		conn:   conn,
		client: proto.NewTychoIndexerClient(conn),
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// GetStatus returns the current status of the Tycho node
func (c *Client) GetStatus(ctx context.Context) (*proto.GetStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	return c.client.GetStatus(ctx, &proto.GetStatusRequest{})
}

// GetLibraryCell fetches a library cell by hash
func (c *Client) GetLibraryCell(ctx context.Context, hash []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	req := &proto.GetLibraryCellRequest{
		Hash: hash,
	}

	resp, err := c.client.GetLibraryCell(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get library cell: %w", err)
	}

	switch result := resp.Library.(type) {
	case *proto.GetLibraryCellResponse_NotFound:
		return nil, fmt.Errorf("library cell not found")
	case *proto.GetLibraryCellResponse_Found:
		return result.Found.Cell, nil
	default:
		return nil, fmt.Errorf("unexpected response type")
	}
}

// GetRawBlockData fetches raw BOC data for debugging purposes
func (c *Client) GetRawBlockData(ctx context.Context, workchain int32, shard uint64, seqno uint32) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	req := &proto.GetBlockRequest{
		Query: &proto.GetBlockRequest_BySeqno{
			BySeqno: &proto.BlockBySeqno{
				Workchain: workchain,
				Shard:     shard,
				Seqno:     seqno,
			},
		},
	}

	stream, err := c.client.GetBlock(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start block stream: %w", err)
	}

	return c.readBlockFromStream(stream)
}

// readBlockFromStream reads block data from the gRPC stream
func (c *Client) readBlockFromStream(stream proto.TychoIndexer_GetBlockClient) ([]byte, error) {
	var totalData []byte

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("stream error: %w", err)
		}

		switch msg := resp.Msg.(type) {
		case *proto.GetBlockResponse_NotFound:
			return nil, BlockNotFoundErr
		case *proto.GetBlockResponse_Found:
			// First chunk with metadata
			if msg.Found.FirstChunk != nil {
				totalData = append(totalData, msg.Found.FirstChunk.Data...)
			}
		case *proto.GetBlockResponse_Chunk:
			// Subsequent chunks
			totalData = append(totalData, msg.Chunk.Data...)
		default:
			return nil, fmt.Errorf("unexpected response type: %T", msg)
		}
	}

	if len(totalData) == 0 {
		return nil, fmt.Errorf("no block data received")
	}

	return totalData, nil
}
