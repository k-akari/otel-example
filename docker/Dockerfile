FROM golang:1.22.3-bookworm as builder
ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/k-akari/otel-example
COPY . .
RUN apt-get update && apt-get install --no-install-recommends -y make
RUN make build

FROM gcr.io/distroless/static-debian12:latest as http_server
COPY --from=builder --chown=nonroot:nonroot /go/src/github.com/k-akari/otel-example/dist/http_server /http_server
USER nonroot
CMD ["/http_server"]

FROM gcr.io/distroless/static-debian12:latest as grpc_server
COPY --from=builder --chown=nonroot:nonroot /go/src/github.com/k-akari/otel-example/dist/grpc_server /grpc_server
USER nonroot
CMD ["/grpc_server"]
