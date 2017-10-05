package main

import (
	"errors"
	"net/url"
	"regexp"
)

var (
	// ErrMalformedURL @TODO
	ErrMalformedURL = errors.New("malformed URL")
	// ErrNotStream @TODO
	ErrNotStream = errors.New("not a stream item")

	urlMatcher    = regexp.MustCompile(`(\w+)\/(\w+)\/(rss\.xml|\d+\.mp3)`)
	streamMatcher = regexp.MustCompile(`^(\d+).mp3$`)
)

// SoundcloudRSSRequest @TODO
type SoundcloudRSSRequest struct {
	Username string
	Playlist string
	Item     string
}

// ParseURL @TODO
func (r *SoundcloudRSSRequest) ParseURL(u *url.URL) error {
	match := urlMatcher.FindStringSubmatch(u.Path)
	if len(match) != 4 {
		return ErrMalformedURL
	}
	r.Username = match[1]
	r.Playlist = match[2]
	r.Item = match[3]
	return nil
}

// IsFeed @TODO
func (r *SoundcloudRSSRequest) IsFeed() bool {
	return "rss.xml" == r.Item
}

// IsStream @TODO
func (r *SoundcloudRSSRequest) IsStream() bool {
	match := streamMatcher.FindStringSubmatch(r.Item)
	return len(match) == 2
}

// StreamID @TODO
func (r *SoundcloudRSSRequest) StreamID() string {
	match := streamMatcher.FindStringSubmatch(r.Item)
	if len(match) != 2 {
		return ""
	}
	return match[1]
}
