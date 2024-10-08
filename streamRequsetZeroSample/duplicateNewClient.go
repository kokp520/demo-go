package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"demo-go/streamRequsetZeroSample/types/myapp/myservice"

	"google.golang.org/grpc"
)

func createClientAndSendRequests(index int) {
	// 建立 gRPC 連線
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 建立 MyServiceClient
	client := myservice.NewMyServiceClient(conn)
	stream, err := client.StreamRequests(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; i < 2; i++ {
			err := stream.Send(&myservice.MyRequest{Data: fmt.Sprintf("Client #%d Request #%d", index, i)})
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
		fmt.Printf("Client #%d Received: %s\n", index, res.Message)
	}
}

func main() {
	// 重複建立多個 MyServiceClient
	for clientIndex := 0; clientIndex < 5; clientIndex++ {
		// 每個 client 獨立的 goroutine
		go createClientAndSendRequests(clientIndex)
		// time.Sleep(1 * time.Second)
	}

	// 為了防止主程式提前結束，等待所有 goroutine 完成
	time.Sleep(20 * time.Second)
}
