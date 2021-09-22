package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	trippb "learn_grpc/01/gen/go"
)

func main() {
	trip := trippb.Trip{
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
	}

	fmt.Println(&trip)
	b, err := proto.Marshal(&trip)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%X\n", b)

	var trip2 trippb.Trip
	err = proto.Unmarshal(b, &trip2)

	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", &trip2)
}
