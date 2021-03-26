package utils

// utility function to initialitze the connection to the server from the client

import (
    "context"
    "log"
	"time"

	"google.golang.org/grpc"
	pb "num/arexservices"

)

const (
	address     = "localhost:50051"
)

// Returns all the necessari objects to connect to the server.
func InitClient() (*grpc.ClientConn,pb.ArexServicesClient,context.Context,context.CancelFunc) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    c := pb.NewArexServicesClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return conn,c,ctx,cancel
}
