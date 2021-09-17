package hackernews

import (
	"context"
	"io"

	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Client) FetchAll(ctx context.Context) ([]Item, error) {
	clientStream, err := c.client.ListAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "streaming all items")
	}

	return collectStreamItems(ctx, clientStream)
}

func (c *Client) FetchStories(ctx context.Context) ([]Item, error) {
	clientStream, err := c.client.ListStories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "streaming story items")
	}

	return collectStreamItems(ctx, clientStream)
}

func (c *Client) FetchJobs(ctx context.Context) ([]Item, error) {
	clientStream, err := c.client.ListJobs(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "streaming job items")
	}

	return collectStreamItems(ctx, clientStream)
}

type itemStream interface {
	Recv() (*pb.Item, error)
}

func collectStreamItems(ctx context.Context, s itemStream) ([]Item, error) {
	items := []Item{}
	for {
		select {
		case <-ctx.Done():
			return items, nil
		default:
			item, err := s.Recv()
			if err != nil {
				if err == io.EOF {
					return items, nil
				}

				return nil, errors.Wrap(err, "receiving items from server")
			}

			items = append(items, ptoh(item))
		}
	}
}
