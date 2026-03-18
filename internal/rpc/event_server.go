package rpc

import (
	"context"

	pb "myCalendar/grpc/pb"
	"myCalendar/internal/rpc/create_event"
	"myCalendar/internal/rpc/delete_event"
	"myCalendar/internal/rpc/get_event"
	"myCalendar/internal/rpc/get_user_events"
	"myCalendar/internal/rpc/update_event"
)

// EventServer — реализует pb.EventServiceServer
// собирает все event ручки в одно место
type EventServer struct {
	pb.UnimplementedEventServiceServer
	createEvent   *create_event.Handler
	getEvent      *get_event.Handler
	getUserEvents *get_user_events.Handler
	updateEvent   *update_event.Handler
	deleteEvent   *delete_event.Handler
}

func NewEventServer(
	createEvent *create_event.Handler,
	getEvent *get_event.Handler,
	getUserEvents *get_user_events.Handler,
	updateEvent *update_event.Handler,
	deleteEvent *delete_event.Handler,
) *EventServer {
	return &EventServer{
		createEvent:   createEvent,
		getEvent:      getEvent,
		getUserEvents: getUserEvents,
		updateEvent:   updateEvent,
		deleteEvent:   deleteEvent,
	}
}

func (s *EventServer) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	return s.createEvent.Handle(ctx, req)
}

func (s *EventServer) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	return s.getEvent.Handle(ctx, req)
}

func (s *EventServer) GetUserEvents(ctx context.Context, req *pb.GetUserEventsRequest) (*pb.GetUserEventsResponse, error) {
	return s.getUserEvents.Handle(ctx, req)
}

func (s *EventServer) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	return s.updateEvent.Handle(ctx, req)
}

func (s *EventServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	return s.deleteEvent.Handle(ctx, req)
}
