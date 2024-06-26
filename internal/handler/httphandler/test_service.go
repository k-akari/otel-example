package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"github.com/k-akari/otel-example/internal/infra/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestService struct {
	db  *database.Client
	tsc testpb.TestServiceClient
}

func NewTestService(db *database.Client, tsc testpb.TestServiceClient) *TestService {
	return &TestService{
		db:  db,
		tsc: tsc,
	}
}

func (h *TestService) PingGRPC(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var b struct {
		Value             string `json:"value"`
		SleepTimeMs       int32  `json:"sleep_time_ms"`
		ErrorCodeReturned int32  `json:"error_code_returned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		respondJSON(ctx, w, &errResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	if err := h.db.Query(ctx); err != nil {
		respondJSON(ctx, w, &errResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	_, err := h.tsc.Ping(ctx, &testpb.PingRequest{
		Value:             b.Value,
		SleepTimeMs:       b.SleepTimeMs,
		ErrorCodeReturned: uint32(b.ErrorCodeReturned),
	})
	if err != nil {
		if status.Code(err) != codes.Code(b.ErrorCodeReturned) {
			respondJSON(ctx, w, &errResponse{Message: err.Error()}, http.StatusInternalServerError)
			return
		} else {
			respondJSON(ctx, w, &errResponse{Message: err.Error()}, http.StatusOK)
			return
		}
	}

	respondJSON(ctx, w, nil, http.StatusOK)
}
