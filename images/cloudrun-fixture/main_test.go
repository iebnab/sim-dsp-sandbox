package main

import (
	"net/http/httptest"
	"testing"
)

func TestEveryRequestReturnsConstantOK(t *testing.T) {
	for _, target := range []string{"/", "/health", "/unmatched?query=value"} {
		for _, method := range []string{"GET", "POST", "DELETE"} {
			request := httptest.NewRequest(method, target, nil)
			response := httptest.NewRecorder()
			serveOK(response, request)
			if response.Code != 200 || response.Body.String() != "ok\n" {
				t.Fatalf("%s %s: status=%d body=%q", method, target, response.Code, response.Body.String())
			}
		}
	}
}

func TestJobModeDoesNotStartListener(t *testing.T) {
	if err := run([]string{"--job"}, "invalid-port"); err != nil {
		t.Fatalf("job must exit successfully without a listener: %v", err)
	}
}

func TestServiceRejectsInvalidPort(t *testing.T) {
	if err := run(nil, "invalid-port"); err == nil {
		t.Fatal("service accepted invalid port")
	}
}
