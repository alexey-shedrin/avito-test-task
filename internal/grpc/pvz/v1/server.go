package pvzv1

import (
	"fmt"
	"log"
	"net"

	"github.com/alexey-shedrin/avito-test-task/internal/handler"
	"google.golang.org/grpc"
)

type PVZServer struct {
	UnimplementedPVZServiceServer
	pvzService handler.PvzService
}

func Start(grpcServerPort string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcServerPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	RegisterPVZServiceServer(server, &PVZServer{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
