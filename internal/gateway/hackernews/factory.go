package hackernews

import (
	"context"

	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Client interface {
	FetchAll(ctx context.Context) ([]models.Item, error)
	FetchStories(ctx context.Context) ([]models.Item, error)
	FetchJobs(ctx context.Context) ([]models.Item, error)
}

type client struct {
	client pb.APIClient
	conn   *grpc.ClientConn
}

func New(addr string) (*client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to grpc server")
	}

	c := pb.NewAPIClient(conn)

	return &client{
		client: c,
		conn:   conn,
	}, nil
}

func (c *client) Close() {
	c.conn.Close()
}
