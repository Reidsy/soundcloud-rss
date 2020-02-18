package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var (
	// ErrMalformedFeedURL is returned when a feed url cannot be used to locate a user and a podcast
	ErrMalformedFeedURL = errors.New("malformed Feed URL")
)

// FeedServerSource describes the functions required to fetch a playlist
type FeedServerSource interface {
	Playlist(username string, playlistName string) (Playlist, error)
}

// FeedServer is a http server that encodes a Playlist from a FeedServerSource
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
		log.Printf("Playlist Error: %s", playlistErr)
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
