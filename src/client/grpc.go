package client

import (
	"context"
	"fmt"

	"github.com/vdbulcke/confkv/src/logger"
	"github.com/vdbulcke/confkv/src/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ConfKVClient
}

func NewClient(addr string) (*Client, error) {
	// opts := []grpc.DialOption{}{
	// 	grpc.WithTransportCredentials(insecure.NewCredentials(),
	// // }

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewConfKVClient(conn)

	c := &Client{
		conn:   conn,
		client: client,
	}

	return c, nil
}

func (c *Client) Get(bkt, key string) ([]byte, error) {

	traceID := logger.NewTraceID()

	md := metadata.New(map[string]string{
		"x-trace-id": fmt.Sprintf("client-%s", traceID),
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := c.client.Get(ctx, &pb.GetMessage{Bucket: bkt, Key: key})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
