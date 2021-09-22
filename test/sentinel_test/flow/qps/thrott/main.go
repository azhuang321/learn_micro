package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"log"
	"time"
)

const resName = "example-flow-qps-resource"

func main() {
	// We should initialize Sentinel first.
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, //匀速通过
			Threshold:              9,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	ch := make(chan interface{})
	for i := 0; i < 10; i++ {
		e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			fmt.Println("限流")
		} else {
			fmt.Println("通过")
			e.Exit()
		}
		time.Sleep(101 * time.Millisecond)
	}

	<-ch
}
