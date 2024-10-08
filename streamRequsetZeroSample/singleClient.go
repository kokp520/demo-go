package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"demo-go/streamRequsetZeroSample/types/myapp/myservice"

	"google.golang.org/grpc"
)

func createClientAndSendRequest(index int) {
	// 建立 gRPC 連線
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := myservice.NewMyServiceClient(conn)
	sameContext := context.Background()
	stream, err := client.StreamRequests(sameContext)

	for i := 0; i < 2; i++ {
		// req := &myservice.MyRequest{Data: fmt.Sprintf("Client #%d Request #%d", index, i)}
		// res, err := client.Send(customContext, req)
		err := stream.Send(&myservice.MyRequest{Data: fmt.Sprintf("Client #%d Request #%d", index, i)})
		if err != nil {
			panic(err)
		}

		// fmt.Printf("Client #%d Received: %s\n", idex, res.Message)
		// time.Sleep(1 * time.Second)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Client #%d Received: %s\n", index, res.Message)
	}
}

func main() {
	for clientIndex := 0; clientIndex < 2; clientIndex++ {
		go createClientAndSendRequest(clientIndex)
	}

	time.Sleep(20 * time.Second)
}
