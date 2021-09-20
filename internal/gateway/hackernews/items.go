package hackernews

import (
	"context"
	"io"

	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *client) FetchAll(ctx context.Context) ([]models.Item, error) {
	clientStream, err := c.client.ListAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "streaming all items")
	}

	return collectStreamItems(ctx, clientStream)
}

func (c *client) FetchStories(ctx context.Context) ([]models.Item, error) {
	clientStream, err := c.client.ListStories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "streaming story items")
	}

	return collectStreamItems(ctx, clientStream)
}

func (c *client) FetchJobs(ctx context.Context) ([]models.Item, error) {
	clientStream, err := c.client.ListJobs(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "streaming job items")
	}

	return collectStreamItems(ctx, clientStream)
}

type itemStream interface {
	Recv() (*pb.Item, error)
}

func collectStreamItems(ctx context.Context, s itemStream) ([]models.Item, error) {
	items := []models.Item{}
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

			items = append(items, models.Ptoi(item))
		}
	}
}
