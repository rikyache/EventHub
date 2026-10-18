package booking

import (
	pb "EventHub/gen/booking/v1"
	"EventHub/internal/repository"
	"EventHub/internal/service"
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type Server struct {
	pb.UnimplementedBookingServiceServer
	reservationService *service.ReservationService
}

func NewServer(reservationService *service.ReservationService) *Server {
	return &Server{
		reservationService: reservationService,
	}

}

func (s *Server) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.CreateReservationResponse, error) {
	eventID, err := uuid.Parse(req.GetEventId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event_id")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	reservationID, err := s.reservationService.Create(
		ctx,
		eventID,
		userID,
	)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReservationUnavailable):
			return nil, status.Error(
				codes.FailedPrecondition,
				"reservation unavailable",
			)
		default:
			slog.Error("failed to create reservation", "error", err)
			return nil, status.Error(
				codes.Internal,
				"internal server error",
			)
		}
	}

	return &pb.CreateReservationResponse{
		ReservationId: reservationID.String(),
	}, nil
}

func (s *Server) ConfirmReservation(ctx context.Context, req *pb.ConfirmReservationRequest) (*pb.ConfirmReservationResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user_id",
		)
	}
	reservationID, err := uuid.Parse(req.GetReservationId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid reservation_id",
		)
	}

	ticketID, err := s.reservationService.Confirm(ctx, reservationID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReservationNotFound):
			return nil, status.Error(
				codes.NotFound,
				"reservation not found",
			)

		case errors.Is(err, repository.ErrReservationExpired):
			return nil, status.Error(
				codes.FailedPrecondition,
				"reservation expired",
			)

		case errors.Is(err, repository.ErrReservationNotPending):
			return nil, status.Error(
				codes.FailedPrecondition,
				"reservation is not pending",
			)

		case errors.Is(err, repository.ErrInsufficientFunds):
			return nil, status.Error(
				codes.FailedPrecondition,
				"insufficient funds",
			)

		default:
			slog.Error("failed to confirm reservation", "error", err)
			return nil, status.Error(
				codes.Internal,
				"internal server error",
			)
		}
	}

	return &pb.ConfirmReservationResponse{
		TicketId: ticketID.String(),
	}, nil
}

func (s *Server) CancelReservation(ctx context.Context, req *pb.CancelReservationRequest) (*pb.CancelReservationResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	reservationID, err := uuid.Parse(req.GetReservationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reservation_id")
	}

	err = s.reservationService.Cancel(ctx, userID, reservationID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReservationNotFound):
			return nil, status.Error(
				codes.NotFound,
				"reservation not found",
			)

		case errors.Is(err, repository.ErrReservationNotPending):
			return nil, status.Error(
				codes.FailedPrecondition,
				"reservation is not pending",
			)

		case errors.Is(err, repository.ErrReservationUnavailable):
			return nil, status.Error(
				codes.FailedPrecondition,
				"reservation unavailable",
			)

		default:
			slog.Error("failed to cancel reservation", "error", err)
			return nil, status.Error(
				codes.Internal,
				"internal server error",
			)
		}
	}

	return &pb.CancelReservationResponse{}, nil
}
