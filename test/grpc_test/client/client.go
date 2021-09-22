package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"mygrpc/proto"
)

type customCredentials struct{}

func (c customCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": "test1",
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (c customCredentials) RequireTransportSecurity() bool {
	return false
}

func main() {
	//intercepter := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//	startTime := time.Now()
	//	md := metadata.New(map[string]string{
	//		"token":"test1",
	//	})
	//	ctx = metadata.NewOutgoingContext(context.Background(),md)
	//	err := invoker(ctx,method,req,reply,cc,opts...)
	//	fmt.Printf("耗时：%s\n",time.Since(startTime))
	//	return err
	//}
	//opt := grpc.WithUnaryInterceptor(intercepter)
	//conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure(),opt)

	cfg := jaegercfg.Configuration{
		ServiceName: "test_jaeger",
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "172.17.0.1:6831",
		},
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(customCredentials{}),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	md := metadata.New(map[string]string{
		"username": "test1",
		"password": "test2",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	r, err := c.SayHello(ctx, &proto.HelloRequest{
		Name:  "test",
		Hobby: []string{"swimming", "running"},
		Sex:   proto.Sex_FEMALE,
		Mp: map[string]string{
			"name": "test1",
			"sex":  "test2",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
	fmt.Println(r.Hobby)
	fmt.Println(r.Sex)
	fmt.Println(r.Mp)
}
