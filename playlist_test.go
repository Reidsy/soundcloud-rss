package main

import (
	"bufio"
	"bytes"
	"io"
	"testing"
	"time"
)

type StubPlaylist struct {
}

func (p *StubPlaylist) Title() string {
	return "My Example Playlist"
}

func (p *StubPlaylist) Description() string {
	return "A description of My Example Playlist"
}

func (p *StubPlaylist) Link() string {
	return "http://example.com/source-of-this-playlist"
}

func (p *StubPlaylist) Author() string {
	return "A cool guy"
}

func (p *StubPlaylist) PubDate() *time.Time {
	publish, _ := time.Parse(time.RFC3339, "2017-10-24T15:04:05Z")
	return &publish
}

func (p *StubPlaylist) LastBuild() *time.Time {
	build, _ := time.Parse(time.RFC3339, "2017-10-24T17:00:25Z")
	return &build
}

func (p *StubPlaylist) Tracks() []Track {
	goodTrack := StubTrack{}
	invalidTrack := StubInvalidTrack{}
	tracks := []Track{}
	tracks = append(tracks, &goodTrack, &invalidTrack)
	return tracks
}

type StubTrack struct {
}

func (t *StubTrack) ID() uint {
	return 12
}

func (t *StubTrack) Title() string {
	return "Track Title"
}

func (t *StubTrack) Description() string {
	return "Track Description"
}

func (t *StubTrack) PubDate() *time.Time {
	publish, _ := time.Parse(time.RFC3339, "2017-10-24T07:54:45Z")
	return &publish
}

func (t *StubTrack) Link() string {
	return "http://example.com/source-of-this-track.html"
}

func (t *StubTrack) Stream() string {
	return "http://example.com/source-of-this-track.mp3"
}

type StubInvalidTrack struct {
	StubTrack
}

func (t *StubInvalidTrack) Description() string {
	return ""
}

type StubPlaylistEncoder struct {
}

func (e *StubPlaylistEncoder) Encode(w io.Writer, p Playlist) error {
	w.Write([]byte(p.Title()))
	return nil
}

const EncodedPodcastPlaylist = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <title>My Example Playlist</title>
    <link>http://example.com/source-of-this-playlist</link>
    <description>A description of My Example Playlist</description>
    <generator>go podcast v1.3.1 (github.com/eduncan911/podcast)</generator>
    <language>en-us</language>
    <lastBuildDate>Tue, 24 Oct 2017 17:00:25 +0000</lastBuildDate>
    <pubDate>Tue, 24 Oct 2017 15:04:05 +0000</pubDate>
    <item>
      <guid>http://example.com/source-of-this-track.mp3</guid>
      <title>Track Title</title>
      <link>http://example.com/source-of-this-track.html</link>
      <description>Track Description</description>
      <pubDate>Tue, 24 Oct 2017 07:54:45 +0000</pubDate>
      <enclosure url="http://example.com/source-of-this-track.mp3" length="0" type="audio/mpeg"></enclosure>
    </item>
  </channel>
</rss>`

func TestEncodePlaylistAsPodcast(t *testing.T) {
	encoder := PodcastPlaylistEncoder{}
	p := StubPlaylist{}
	receivedData := bytes.Buffer{}
	w := bufio.NewWriter(&receivedData)
	encoder.Encode(w, &p)

	if receivedData.String() != EncodedPodcastPlaylist {
		t.Fatalf("Unexpected data received. Got: %s", receivedData.String())
	}
}
