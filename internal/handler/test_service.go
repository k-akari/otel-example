package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestService struct {
	testpb.UnimplementedTestServiceServer
}

func NewTestService() *TestService {
	return &TestService{}
}

func (h *TestService) Ping(_ context.Context, in *testpb.PingRequest) (*testpb.PingResponse, error) {
	if in.Value == "panic" {
		panic("something went wrong")
	}

	if in.SleepTimeMs > 0 {
		time.Sleep(time.Duration(in.SleepTimeMs) * time.Millisecond)
	}

	if in.ErrorCodeReturned != 0 {
		return nil, status.Error(codes.Code(in.ErrorCodeReturned), fmt.Sprintf("something went wrong: %s", in.Value))
	}

	return nil, nil
}

func (h *TestService) PingEmpty(_ context.Context, _ *testpb.PingEmptyRequest) (*testpb.PingEmptyResponse, error) {
	return nil, nil
}

func (h *TestService) PingError(_ context.Context, in *testpb.PingErrorRequest) (*testpb.PingErrorResponse, error) {
	return nil, errors.New("something went wrong")
}

func (h *TestService) PingList(_ *testpb.PingListRequest, _ testpb.TestService_PingListServer) error {
	return nil
}

func (h *TestService) PingStream(stream testpb.TestService_PingStreamServer) error {
	return nil
}
