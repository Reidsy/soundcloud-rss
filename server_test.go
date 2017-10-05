package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerMethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest("POST", "/reidsy/likes/rss.xml", strings.NewReader(`{}`))
	res := httptest.NewRecorder()
	s := SoundcloudRSSServer{}
	s.ServeHTTP(res, req)

	if res.Code != 405 {
		t.Fatalf("Expected http status 405, Got %d", res.Code)
	}
}

func TestServerMalformedURL(t *testing.T) {
	req, _ := http.NewRequest("GET", "/invalid/url.html", nil)
	res := httptest.NewRecorder()
	s := SoundcloudRSSServer{}
	s.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Fatalf("Expected http status 404, Got %d", res.Code)
	}
}

func TestServerFeed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reidsy/likes/rss.xml", nil)
	res := httptest.NewRecorder()
	s := SoundcloudRSSServer{}
	s.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Fatalf("Expected http status 200, Got %d", res.Code)
	}
}

func TestServerStream(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reidsy/likes/123.mp3", nil)
	res := httptest.NewRecorder()
	s := SoundcloudRSSServer{}
	s.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Fatalf("Expected http status 200, Got %d", res.Code)
	}
}
