package interceptor

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func UnaryClientInterceptors(tr trace.Tracer) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx, span := tr.Start(ctx, method)
			defer span.End()
			return invoker(ctx, method, req, reply, cc, opts...)
		},
	)
}
