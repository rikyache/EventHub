package booking

import (
	pb "EventHub/gen/booking/v1"
	"EventHub/internal/service"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func Run(address string, reservationService *service.ReservationService) error {
	// открывает TCP порт
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}
	defer lis.Close()

	// создание grpc
	grpcServer := grpc.NewServer()

	bookingServer := NewServer(reservationService)

	//передача при запросах bookingServer использовтаь сервер
	pb.RegisterBookingServiceServer(
		grpcServer,
		bookingServer,
	)

	return grpcServer.Serve(lis)
}
