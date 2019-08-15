package main

import (
	"io"
	"log"
	"time"

	"github.com/eduncan911/podcast"
)

type Playlist interface {
	Title() string
	Description() string
	Link() string
	Author() string
	PubDate() *time.Time
	LastBuild() *time.Time
	Tracks() []Track
}

type Track interface {
	ID() uint
	Title() string
	Description() string
	PubDate() *time.Time
	Link() string
	Image() string
	Stream() string
}

type PlaylistEncoder interface {
	Encode(io.Writer, Playlist) error
}

type PodcastPlaylistEncoder struct {
}

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
