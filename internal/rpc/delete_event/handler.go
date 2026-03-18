package delete_event

import (
	"context"

	"go.uber.org/zap"
	pb "myCalendar/grpc/pb"
	"myCalendar/internal/apperrors"
	"myCalendar/internal/ctxutil"
)

type Handler struct {
	repo repo
	log  *zap.Logger
}

func New(repo repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) Handle(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
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

	if err := h.repo.Delete(ctx, req.Id); err != nil {
		return nil, err
	}

	h.log.Info("event deleted", zap.String("id", req.Id))

	return &pb.DeleteEventResponse{Success: true}, nil
}
