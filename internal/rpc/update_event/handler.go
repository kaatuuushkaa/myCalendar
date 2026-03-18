package update_event

import (
	"context"
	"myCalendar/internal/domain"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/ctxutil"
)

const timeLayout = time.RFC3339

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) Handle(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	if err := validate(req); err != nil {
		return nil, err
	}

	userID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	e, err := h.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if e.UserID != userID {
		return nil, apperrors.ErrAccessDenied
	}

	startAt, err := time.Parse(timeLayout, req.StartAt)
	if err != nil {
		return nil, apperrors.ErrInvalidTime
	}

	endAt, err := time.Parse(timeLayout, req.EndAt)
	if err != nil {
		return nil, apperrors.ErrInvalidTime
	}

	if endAt.Before(startAt) {
		return nil, apperrors.ErrEndBeforeStart
	}

	e.Title = req.Title
	e.Description = req.Description
	e.StartAt = startAt
	e.EndAt = endAt
	e.EventDate = startAt.Format("2006-01-02")

	event_updated, err := h.repo.Update(ctx, e.ID, e.Title, e.UserID, e.EventDate, e.StartAt, e.EndAt)
	if err != nil {
		return nil, err
	}

	h.log.Info("event updated", zap.String("id", event_updated.ID))

	return &pb.UpdateEventResponse{Success: true, Event: toProto(event_updated)}, nil
}

func validate(req *pb.UpdateEventRequest) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "event id is required")
	}
	if req.Title == "" {
		return status.Error(codes.InvalidArgument, "title is required")
	}
	if req.StartAt == "" {
		return status.Error(codes.InvalidArgument, "start_at is required")
	}
	if req.EndAt == "" {
		return status.Error(codes.InvalidArgument, "end_at is required")
	}
	return nil
}

func toProto(e domain.Event) *pb.Event {
	return &pb.Event{
		Id:          e.ID,
		UserId:      e.UserID,
		Title:       e.Title,
		Description: e.Description,
		StartAt:     e.StartAt.Format(timeLayout),
		EndAt:       e.EndAt.Format(timeLayout),
		EventDate:   e.EventDate,
	}
}
