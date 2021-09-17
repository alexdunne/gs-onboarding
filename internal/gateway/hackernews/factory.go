package hackernews

import (
	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Client struct {
	client pb.APIClient
	conn   *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to grpc server")
	}

	client := pb.NewAPIClient(conn)

	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
