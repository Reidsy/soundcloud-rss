package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var (
	ErrMalformedFeedURL = errors.New("malformed Feed URL")
)

type FeedServerSource interface {
	Playlist(username string, playlistName string) (Playlist, error)
}
type FeedServer struct {
	Source  FeedServerSource
	Encoder PlaylistEncoder
}

func (f FeedServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		f.sendStatus(w, http.StatusMethodNotAllowed)
		return
	}

	username, playlistName, urlError := f.parseURL(r.URL)
	if urlError == ErrMalformedFeedURL {
		f.sendStatus(w, http.StatusNotFound)
		return
	}

	playlist, playlistErr := f.Source.Playlist(username, playlistName)
	if playlistErr != nil {
		f.sendStatus(w, http.StatusServiceUnavailable)
		return
	}
	
	f.Encoder.Encode(w, playlist)
}

func (f FeedServer) sendStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s\n", http.StatusText(status))
}

func (f FeedServer) parseURL(u *url.URL) (string, string, error) {
	matcher := regexp.MustCompile(`\/(\w+)\/(\w+)\/rss\.xml$`)
	match := matcher.FindStringSubmatch(u.Path)
	if len(match) != 3 {
		return "", "", ErrMalformedFeedURL
	}
	return match[1], match[2], nil
}
