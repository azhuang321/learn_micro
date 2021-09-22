package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	trippb "learn_grpc/03/gen/go"
	"log"
	"net"
	"net/http"
	"time"
)

type TripService struct{}

func (t *TripService) GetTrip(ctx context.Context, request *trippb.GetTripRequest) (*trippb.GetTripResponse, error) {
	return &trippb.GetTripResponse{
		Id: request.Id,
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
	myChan := make(chan interface{})
	go func() {
		go startGRPCGateway()
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic(err)
		}
		g := grpc.NewServer()
		tripSrv := &TripService{}
		trippb.RegisterTripServiceServer(g, tripSrv)
		log.Fatal(g.Serve(listener))
	}()

	time.Sleep(time.Microsecond * 500)

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
	<-myChan
}

// curl http://127.0.0.1:8082/trip/1111 测试
func startGRPCGateway() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			EnumsAsInts: true,
			OrigName:    true,
		},
	))
	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		":8081",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(":8082", mux)
	if err != nil {
		panic(err)
	}
}
