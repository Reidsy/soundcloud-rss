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
	User         soundcloudUser
	PlaylistName string
	tracks       []Track
}

func (p SoundcloudPlaylist) Title() string {
	return strings.Title(p.PlaylistName)
}

func (p SoundcloudPlaylist) Link() string {
	return fmt.Sprintf("https://soundcloud.com/%s/likes", p.User.Permalink)
}

func (p SoundcloudPlaylist) Description() string {
	return p.Title()
}

func (p SoundcloudPlaylist) Author() string {
	return strings.Title(p.User.Username)
}

func (p SoundcloudPlaylist) PubDate() *time.Time {
	t := time.Now()
	return &t
}

func (p SoundcloudPlaylist) LastBuild() *time.Time {
	t := time.Now()
	return &t
}

func (p SoundcloudPlaylist) Tracks() []Track {
	return p.tracks
}

func (s *SoundcloudSource) Playlist(username string, playlistName string) (Playlist, error) {
	user, err := s.fetchUser(username)
	if err != nil {
		return SoundcloudPlaylist{}, err
	}
	
	likes, err := s.fetchUserLikes(user)
	if err != nil {
		return SoundcloudPlaylist{}, err
	}
	
	playlist := SoundcloudPlaylist{User: user, PlaylistName: playlistName, tracks: likes}
	return playlist, nil
}

type soundcloudUser struct {
	ID                   uint   `json:"id"`
	Permalink            string `json:"permalink"`
	Username             string `json:"username"`
	URI                  string `json:"uri"`
	PublicFavoritesCount uint   `json:"public_favorites_count"`
}

func (s *SoundcloudSource) fetchUser(username string) (soundcloudUser, error) {
	u, _ := url.Parse("https://api.soundcloud.com/resolve")
	q := u.Query()
	q.Set("client_id", s.ClientID)
	q.Set("url", fmt.Sprintf("https://soundcloud.com/%s", username))
	u.RawQuery = q.Encode()

	user := soundcloudUser{}
	
	resp, err := s.Client.Get(u.String())
	if err != nil {
		return user, err
	}
	if resp.StatusCode != 200 {
		return user, fmt.Errorf("fetchUser received %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	return user, err
}

type SoundcloudTrack struct {
	id          uint
	title       string
	description string
	pubdate     time.Time
	link        string
	image       string
	stream      string
}

func (st SoundcloudTrack) ID() uint {
	return st.id
}

func (st SoundcloudTrack) Title() string {
	return st.title
}

func (st SoundcloudTrack) Description() string {
	if st.description != "" {
		return st.description
	} else {
		return st.title
	}
}

func (st SoundcloudTrack) PubDate() *time.Time {
	return &st.pubdate
}

func (st SoundcloudTrack) Link() string {
	return st.link
}

func (st SoundcloudTrack) Image() string {
	return st.image
}

func (st SoundcloudTrack) Stream() string {
	return st.stream
}

type soundcloudUserLikes struct {
	Collection []struct {
		CreatedAt time.Time `json:"created_at"`
		Track     struct {
			Id           uint   `json:"id"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			PermalinkUrl string `json:"permalink_url"`
			ArtworkUrl   string `json:"artwork_url"`
		} `json:"track"`
	} `json:"collection"`
	NextHref string `json:"next_href"`
}

func (s *SoundcloudSource) fetchUserLikes(user soundcloudUser) ([]Track, error) {
	// @TODO: there are "playlist" entries in the likes endpoint. These will be skipped during encoding
	u, _ := url.Parse(fmt.Sprintf("https://api-v2.soundcloud.com/users/%d/likes", user.ID))

	q := u.Query()
	q.Set("offset", fmt.Sprint(0))
	q.Set("limit", fmt.Sprint(user.PublicFavoritesCount))
	q.Set("client_id", s.ClientID)
	u.RawQuery = q.Encode()

	tracks := []Track{}
	likes := soundcloudUserLikes{}

	resp, err := s.Client.Get(u.String())
	if err != nil {
		return tracks, err
	}
	if resp.StatusCode != 200 {
		return tracks, fmt.Errorf("fetchUserLikes received %s", resp.Status)
	}
	
	json.NewDecoder(resp.Body).Decode(&likes)
	for _, like := range likes.Collection {
		track := SoundcloudTrack{
			id:          like.Track.Id,
			title:       like.Track.Title,
			description: like.Track.Description,
			pubdate:     like.CreatedAt,
			link:        like.Track.PermalinkUrl,
			image:       like.Track.ArtworkUrl,
			stream:      fmt.Sprintf("%s/%d.mp3", s.MediaSource, like.Track.Id),
		}
		tracks = append(tracks, &track)
	}

	return tracks, nil
}

type soundcloudStreams struct {
	URL string `json:"http_mp3_128_url"`
}

func (s *SoundcloudSource) StreamURL(streamID string) string {
	streamsURL := fmt.Sprintf("https://api.soundcloud.com/tracks/%s/streams?client_id=%s", streamID, s.ClientID)
	resp, _ := s.Client.Get(streamsURL)
	streams := soundcloudStreams{}
	json.NewDecoder(resp.Body).Decode(&streams)
	return streams.URL
}
