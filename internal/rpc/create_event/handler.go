package create_event

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/ctxutil"
	"myCalendar/internal/domain"
	"time"
)

const timeLayout = time.RFC3339 // "2026-03-20T14:00:00Z"

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) Handle(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	if err := validate(req); err != nil {
		return nil, err
	}

	userID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	startAt, err := time.Parse(timeLayout, req.StartAt)
	if err != nil {
		return nil, apperrors.ErrInvalidTime
	}

	endAt, err := time.Parse(timeLayout, req.EndAt)
	if err != nil {
		return nil, apperrors.ErrInvalidTime
	}

	e := domain.Event{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		StartAt:     startAt,
		EndAt:       endAt,
		EventDate:   startAt.Format("2006-01-02"),
	}

	if err = h.repo.Create(ctx, e); err != nil {
		h.log.Warn("failed to create event", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	h.log.Info("event_created", zap.String("id", e.ID), zap.String("user_id", userID))

	return &pb.CreateEventResponse{Success: true, Id: e.ID}, nil
}

func validate(req *pb.CreateEventRequest) error {
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
