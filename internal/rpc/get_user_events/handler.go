package get_user_events

import (
	"context"
	"myCalendar/internal/domain"
	"time"

	"go.uber.org/zap"

	pb "myCalendar/grpc/pb"
	"myCalendar/internal/ctxutil"
)

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) Handle(ctx context.Context, _ *pb.GetUserEventsRequest) (*pb.GetUserEventsResponse, error) {
	userID, err := ctxutil.UserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	events, err := h.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// конвертируем каждый доменный ивент в proto
	protoEvents := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		protoEvents = append(protoEvents, toProto(e))
	}

	return &pb.GetUserEventsResponse{Events: protoEvents}, nil
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

const timeLayout = time.RFC3339
