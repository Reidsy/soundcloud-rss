package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SoundcloudSource struct {
	MediaSource  string
	Client       http.Client
	ClientID     string
	ClientSecret string
}

type SoundcloudPlaylist struct {
	User         SoundcloudUser
	PlaylistName string
	tracks       []Track
}

func (p *SoundcloudPlaylist) Title() string {
	return strings.Title(p.PlaylistName)
}

func (p *SoundcloudPlaylist) Link() string {
	return fmt.Sprintf("https://soundcloud.com/%s/likes", p.User.Permalink)
}

func (p *SoundcloudPlaylist) Description() string {
	return p.Title()
}

func (p *SoundcloudPlaylist) Author() string {
	return strings.Title(p.User.Username)
}

func (p *SoundcloudPlaylist) PubDate() *time.Time {
	t := time.Now()
	return &t
}

func (p *SoundcloudPlaylist) LastBuild() *time.Time {
	t := time.Now()
	return &t
}

func (p *SoundcloudPlaylist) Tracks() []Track {
	return p.tracks
}

func (s *SoundcloudSource) Playlist(username string, playlistName string) Playlist {
	user := s.User(username)
	playlist := SoundcloudPlaylist{User: user, PlaylistName: playlistName}
	playlist.tracks = s.likes(user)

	return &playlist
}

type SoundcloudUser struct {
	ID                   uint   `json:"id"`
	Permalink            string `json:"permalink"`
	Username             string `json:"username"`
	URI                  string `json:"uri"`
	PublicFavoritesCount uint   `json:"public_favorites_count"`
}

func (s *SoundcloudSource) User(username string) SoundcloudUser {
	u, _ := url.Parse("https://api.soundcloud.com/resolve")
	q := u.Query()
	q.Set("client_id", s.ClientID)
	q.Set("url", fmt.Sprintf("https://soundcloud.com/%s", username))
	u.RawQuery = q.Encode()

	resp, _ := s.Client.Get(u.String())
	user := SoundcloudUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	return user
}

type SoundcloudAPILikeRequest struct {
	Collection []SoundcloudAPILike `json:"collection"`
	NextHref   string              `json:"next_href"`
}

type SoundcloudAPILike struct {
	CreatedAt time.Time       `json:"created_at"`
	Track     SoundcloudTrack `json:"track"`
}

type SoundcloudTrack struct {
	Fid          uint   `json:"id"`
	Ftitle       string `json:"title"`
	Fdescription string `json:"description"`
	pubdate      time.Time
	Flink        string `json:"permalink_url"`
	stream       string
}

func (t *SoundcloudTrack) ID() uint {
	return t.Fid
}

func (t *SoundcloudTrack) Title() string {
	return t.Ftitle
}

func (t *SoundcloudTrack) Description() string {
	if len(t.Fdescription) == 0 {
		return t.Ftitle
	}
	return t.Fdescription
}

func (t *SoundcloudTrack) PubDate() *time.Time {
	return &t.pubdate
}

func (t *SoundcloudTrack) Link() string {
	return t.Flink
}

func (t *SoundcloudTrack) Stream() string {
	return t.stream
}

func (s *SoundcloudSource) likes(user SoundcloudUser) []Track {
	u, _ := url.Parse(fmt.Sprintf("https://api-v2.soundcloud.com/users/%d/likes", user.ID))

	q := u.Query()
	q.Set("offset", fmt.Sprint(0))
	q.Set("limit", fmt.Sprint(user.PublicFavoritesCount))
	q.Set("client_id", s.ClientID)
	u.RawQuery = q.Encode()

	resp, _ := s.Client.Get(u.String())

	likeRequest := SoundcloudAPILikeRequest{}
	json.NewDecoder(resp.Body).Decode(&likeRequest)
	tracks := []Track{}
	startDate := time.Now()
	for idx, like := range likeRequest.Collection {
		track := like.Track
		track.stream = fmt.Sprintf("%s/%d.mp3", s.MediaSource, track.Fid)
		track.pubdate = startDate.AddDate(0, 0, -idx)
		tracks = append(tracks, &track)
	}
	return tracks
}

type SoundcloudStreams struct {
	URL string `json:"http_mp3_128_url"`
}

func (s *SoundcloudSource) StreamURL(streamID string) string {
	streamsURL := fmt.Sprintf("https://api.soundcloud.com/tracks/%s/streams?client_id=%s", streamID, s.ClientID)
	resp, _ := s.Client.Get(streamsURL)
	streams := SoundcloudStreams{}
	json.NewDecoder(resp.Body).Decode(&streams)
	return streams.URL
}
