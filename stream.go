package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var (
	// ErrMalformedStreamURL is returned when a streamID cannot be pulled from the URL
	ErrMalformedStreamURL = errors.New("malformed Stream URL")
)

// StreamServerSource describes the functions required for a client to fetch the podcast track
type StreamServerSource interface {
	StreamURL(streamID string) string
}

// StreamServer is a http server that redirects clients to the url of a stream
type StreamServer struct {
	Source StreamServerSource
}

func (s StreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendStatus(w, http.StatusMethodNotAllowed)
		return
	}

	streamID, urlError := s.parseURL(r.URL)
	if urlError == ErrMalformedStreamURL {
		s.sendStatus(w, http.StatusNotFound)
		return
	}

	s.redirectToStream(w, streamID)
}

func (s StreamServer) sendStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s\n", http.StatusText(status))
}

func (s StreamServer) parseURL(u *url.URL) (string, error) {
	matcher := regexp.MustCompile(`\/(\d+)\.mp3$`)
	match := matcher.FindStringSubmatch(u.Path)
	if len(match) != 2 {
		return "", ErrMalformedStreamURL
	}
	return match[1], nil
}

func (s StreamServer) redirectToStream(w http.ResponseWriter, streamID string) {
	url := s.Source.StreamURL(streamID)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}
