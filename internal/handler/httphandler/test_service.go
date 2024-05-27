package httphandler

import "net/http"

type TestService struct {
}

func NewTestService() *TestService {
	return &TestService{}
}

func (h *TestService) PingGRPC(w http.ResponseWriter, r *http.Request) {
	return
}
