package main

import (
	"io"
	"log"
	"time"

	"github.com/eduncan911/podcast"
)

// Playlist contains information about a collection of Tracks
type Playlist interface {
	Title() string
	Description() string
	Link() string
	Author() string
	PubDate() *time.Time
	LastBuild() *time.Time
	Tracks() []Track
}

// Track contains information about particular song
type Track interface {
	ID() uint
	Title() string
	Description() string
	PubDate() *time.Time
	Link() string
	Image() string
	Stream() string
}

// PlaylistEncoder takes a Playlist and encoodes it any format sending the output to io.Writer
type PlaylistEncoder interface {
	Encode(io.Writer, Playlist) error
}

// PodcastPlaylistEncoder is a particular type of encoder taking a Playlist and converting it to an rss podcast feed
type PodcastPlaylistEncoder struct {
}

// Encode takes a playlist and encodes it in an rss podcast format writing the output to the provided io.Writer
func (e PodcastPlaylistEncoder) Encode(w io.Writer, p Playlist) error {
	cast := podcast.New(p.Title(), p.Link(), p.Description(), p.PubDate(), p.LastBuild())

	for _, t := range p.Tracks() {
		item := podcast.Item{
			Title:       t.Title(),
			Description: t.Description(),
			PubDate:     t.PubDate(),
			Link:        t.Link(),
		}
		item.AddEnclosure(t.Stream(), podcast.MP3, 0)
		item.AddImage(t.Image())
		if idx, err := cast.AddItem(item); err != nil {
			log.Printf("Skipped adding item at index %d. Playlist: %s Err: %s", idx, p.Link(), err)
		}
	}

	return cast.Encode(w)
}
