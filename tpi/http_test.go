package tpi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseRecorder(t *testing.T) {

	w := httptest.NewRecorder()
	rsr := responseStatusRecorder{ResponseWriter: w}
	rsr.Write([]byte("ok"))
	if rsr.status != http.StatusOK {
		t.Error("200 OK not auto generated")
	}

	w = httptest.NewRecorder()
	rsr = responseStatusRecorder{ResponseWriter: w}
	rsr.WriteHeader(http.StatusBadRequest)
	rsr.Write([]byte("not ok"))
	if rsr.status != http.StatusBadRequest {
		t.Error("http bad request not set")
	}
}
