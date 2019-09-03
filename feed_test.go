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
	s := FeedServer{}
	s.ServeHTTP(res, req)

	if res.Code != 405 {
		t.Fatalf("Expected http status 405, Got %d", res.Code)
	}
}

func TestServerMalformedURL(t *testing.T) {
	req, _ := http.NewRequest("GET", "/invalid/url.html", nil)
	res := httptest.NewRecorder()
	s := FeedServer{}
	s.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Fatalf("Expected http status 404, Got %d", res.Code)
	}
}

type StubFeedServerSource struct {
	playlist Playlist
}

func (s *StubFeedServerSource) Playlist(username string, playlistName string) (Playlist, error) {
	return s.playlist, nil
}

func TestServerFeed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reidsy/likes/rss.xml", nil)
	res := httptest.NewRecorder()

	playlist := StubPlaylist{}
	source := StubFeedServerSource{&playlist}
	encoder := StubPlaylistEncoder{}
	s := FeedServer{&source, &encoder}
	s.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Fatalf("Expected http status 200, Got %d", res.Code)
	}

	if res.Body.String() != playlist.Title() {
		t.Fatalf("Unexpected data received. Got: %s", res.Body.String())
	}
}
