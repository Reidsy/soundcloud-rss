package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamServerMethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/1234.mp3", strings.NewReader(`{}`))
	res := httptest.NewRecorder()
	s := StreamServer{}
	s.ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected http status %d. Got %d", http.StatusMethodNotAllowed, res.Code)
	}
}

func TestStreamServerNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/1234.m4a", nil)
	res := httptest.NewRecorder()
	s := StreamServer{}
	s.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("Expected http status %d. Got %d", http.StatusNotFound, res.Code)
	}
}

type StubStreamServerSource struct {
	URL string
}

func (s *StubStreamServerSource) StreamURL(streamID string) string {
	return s.URL
}

func TestStreamServerRedirectToStream(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/1234.mp3", nil)
	res := httptest.NewRecorder()
	source := StubStreamServerSource{"http://example.com/mysong.mp3"}
	s := StreamServer{&source}
	s.ServeHTTP(res, req)

	if res.Code != http.StatusMovedPermanently {
		t.Fatalf("Expected http status %d. Got %d", http.StatusMovedPermanently, res.Code)
	}

	locationHeader := res.Header().Get("Location")
	if locationHeader != source.URL {
		t.Fatalf("Expected header 'Location: %s'. Got 'Location: %s'", source.URL, locationHeader)
	}
}
