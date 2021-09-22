package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	trippb "learn_grpc/01/gen/go"
	"log"
	"net"
	"time"
)

type TripService struct{}

func (t *TripService) GetTrip(ctx context.Context, request *trippb.GetTripRequest) (*trippb.GetTripResponse, error) {
	return &trippb.GetTripResponse{
		Id: "1",
		Trip: &trippb.Trip{
			Start: "abc",
			StartPos: &trippb.Location{
				Latitude:  10000,
				Longitude: 2000,
			},
			End: "def",
			EndPos: &trippb.Location{
				Latitude:  10000,
				Longitude: 2000,
			},
			DurationSec: 3600,
			FeeCent:     10000,
			PathLocations: []*trippb.Location{
				{
					Latitude:  10000,
					Longitude: 2000,
				},
				{
					Latitude:  10000,
					Longitude: 2000,
				},
			},
			Status: trippb.TripStatus_FINISHED,
		},
	}, nil
}

func main() {
	go func() {
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic(err)
		}
		g := grpc.NewServer()
		tripSrv := &TripService{}
		trippb.RegisterTripServiceServer(g, tripSrv)
		log.Fatal(g.Serve(listener))
	}()

	time.Sleep(time.Second)

	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	tripClient := trippb.NewTripServiceClient(conn)

	resp, err := tripClient.GetTrip(context.Background(), &trippb.GetTripRequest{Id: "1"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
