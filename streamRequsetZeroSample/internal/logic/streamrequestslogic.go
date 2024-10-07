package logic

import (
	"context"
	"io"

	"demo-go/streamRequsetZeroSample/internal/svc"
	"demo-go/streamRequsetZeroSample/types/myapp/myservice"

	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
)

type StreamRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStreamRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StreamRequestsLogic {
	return &StreamRequestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StreamRequestsLogic) StreamRequests(stream myservice.MyService_StreamRequestsServer) error {
	trace := trace.SpanContextFromContext(l.ctx).TraceID().String()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// 讀取完畢，結束通訊
			return nil
		}
		if err != nil {
			return err
		}

		// fmt.Printf("Received request: %s\n", req.Data)
		logx.Infof("Received request: %s, trace: %s", req.Data, trace)

		// 回應客戶端
		err = stream.Send(&myservice.MyResponse{
			Message: "Response to " + req.Data,
		})
		if err != nil {
			return err
		}
	}

	// default return
	return nil
}
