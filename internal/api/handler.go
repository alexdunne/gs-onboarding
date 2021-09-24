package api

import (
	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handler contains the endpoint handlers
type Handler struct {
	pb.UnimplementedAPIServer
	Cache Cache
}

// ListAll streams a collection of items to a client
func (h Handler) ListAll(empty *emptypb.Empty, s pb.API_ListAllServer) error {
	items, err := h.Cache.GetAll(s.Context())
	if err != nil {
		return errors.Wrap(err, "fetching all items")
	}

	for _, v := range items {
		if err := s.Send(models.Itop(v)); err != nil {
			return errors.Wrap(err, "streaming item to client")
		}
	}

	return nil
}

// ListStories streams a collection of story items to a client
func (h Handler) ListStories(empty *emptypb.Empty, s pb.API_ListStoriesServer) error {
	items, err := h.Cache.GetStories(s.Context())
	if err != nil {
		return errors.Wrap(err, "fetching all items")
	}

	for _, v := range items {
		if err := s.Send(models.Itop(v)); err != nil {
			return errors.Wrap(err, "streaming item to client")
		}
	}

	return nil
}

// ListJobs streams a collection of job items to a client
func (h Handler) ListJobs(empty *emptypb.Empty, s pb.API_ListJobsServer) error {
	items, err := h.Cache.GetJobs(s.Context())
	if err != nil {
		return errors.Wrap(err, "fetching all items")
	}

	for _, v := range items {
		if err := s.Send(models.Itop(v)); err != nil {
			return errors.Wrap(err, "streaming item to client")
		}
	}

	return nil
}
