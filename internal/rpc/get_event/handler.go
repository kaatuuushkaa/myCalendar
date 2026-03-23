package get_event

import (
	"context"
	"go.uber.org/zap"
	"myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/ctxutil"
	"myCalendar/internal/domain"
	"time"
)

const timeLayout = time.RFC3339

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{
		repo: repo,
		log:  log}
}

func (h *Handler) Handle(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	if req.Id == "" {
		return nil, apperrors.ErrInvadArgument
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

	return &pb.GetEventResponse{Event: toProto(e)}, nil
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
