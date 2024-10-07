package main

import (
	"context"
	"fmt"
	"io"

	// "myapp/myservice"
	"demo-go/streamRequsetZeroSample/types/myapp/myservice"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := myservice.NewMyServiceClient(conn)
	stream, err := client.StreamRequests(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; i < 10; i++ {
			err := stream.Send(&myservice.MyRequest{Data: fmt.Sprintf("Request #%d", i)})
			if err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received: %s\n", res.Message)
	}
}
