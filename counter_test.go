package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCountHandlerOk(t *testing.T) {
	expStCode := 200
	b := strings.NewReader("{\"code\": 401 }")
	req := httptest.NewRequest(http.MethodPut, "/counter", b)
	w := httptest.NewRecorder()
	CountHandler(w, req)
	if w.Code != 200 {
		t.Errorf("expected status code %d\n", expStCode)
	}
}

func TestCountHandlerMethodNotAllowed(t *testing.T) {
	expStCode := 405
	b := strings.NewReader("{\"code\": 401 }")
	req := httptest.NewRequest(http.MethodPost, "/counter", b)
	w := httptest.NewRecorder()
	CountHandler(w, req)
	if w.Code != 200 {
		t.Errorf("expected status code %d\n", expStCode)
	}
}

func TestCountHandlerBadRequest(t *testing.T) {
	expStCode := 400
	b := strings.NewReader("{\"code\": abc }")
	req := httptest.NewRequest(http.MethodPut, "/counter", b)
	w := httptest.NewRecorder()
	CountHandler(w, req)
	if w.Code != expStCode {
		t.Errorf("expected status code %d\n", expStCode)
	}
}