package api

import (
	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
	"github.com/alexdunne/gs-onboarding/internal/database"
	"github.com/alexdunne/gs-onboarding/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedAPIServer
	DB database.Database
}

func (h Handler) ListAll(empty *emptypb.Empty, s pb.API_ListAllServer) error {
	items, err := h.DB.GetAll(s.Context())
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

func (h Handler) ListStories(empty *emptypb.Empty, s pb.API_ListStoriesServer) error {
	items, err := h.DB.GetStories(s.Context())
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

func (h Handler) ListJobs(empty *emptypb.Empty, s pb.API_ListJobsServer) error {
	items, err := h.DB.GetJobs(s.Context())
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
