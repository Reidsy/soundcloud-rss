package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SoundcloudSource is a ServerSource that returns playlists from Soundcloud
type SoundcloudSource struct {
	MediaSource  string
	Client       http.Client
	ClientID     string
	ClientSecret string
}

// SoundcloudPlaylist is a Playlist used to fetch playlists and tracks from Soundcloud
type SoundcloudPlaylist struct {
	User         soundcloudUser
	PlaylistName string
	tracks       []Track
}

// Title is the title of a playlist
func (p SoundcloudPlaylist) Title() string {
	return strings.Title(p.PlaylistName)
}

// Link is the url to a soundcloud playlist
func (p SoundcloudPlaylist) Link() string {
	return fmt.Sprintf("https://soundcloud.com/%s/likes", p.User.Permalink)
}

// Description of the playlist
func (p SoundcloudPlaylist) Description() string {
	return p.Title()
}

// Author of the playlist
func (p SoundcloudPlaylist) Author() string {
	return strings.Title(p.User.Username)
}

// PubDate the date the podcast was last updated. Returns the current time
func (p SoundcloudPlaylist) PubDate() *time.Time {
	t := time.Now()
	return &t
}

// LastBuild the date the podcast was last updated. Returns the current time
func (p SoundcloudPlaylist) LastBuild() *time.Time {
	t := time.Now()
	return &t
}

// Tracks is an array of tracks which will be represented as items in a podcast
func (p SoundcloudPlaylist) Tracks() []Track {
	return p.tracks
}

// Playlist returns the SoundcloudPlaylist object
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

// SoundcloudTrack is an implementation of Track. It represents an individual Soundcloud song
type SoundcloudTrack struct {
	id          uint
	title       string
	description string
	pubdate     time.Time
	link        string
	image       string
	stream      string
}

// ID is Soundclouds internal ID for a song
func (st SoundcloudTrack) ID() uint {
	return st.id
}

// Title is the name of a song
func (st SoundcloudTrack) Title() string {
	return st.title
}

// Description contains more information about a song, this is sometimes blank
func (st SoundcloudTrack) Description() string {
	if st.description != "" {
		return st.description
	}
	return st.title
}

// PubDate is when the song was added or liked
func (st SoundcloudTrack) PubDate() *time.Time {
	return &st.pubdate
}

// Link to the song on Soundcloud's website
func (st SoundcloudTrack) Link() string {
	return st.link
}

// Image contains artwork for the song
func (st SoundcloudTrack) Image() string {
	return st.image
}

// Stream contains the URL so a podcast client can download the song
func (st SoundcloudTrack) Stream() string {
	return st.stream
}

type soundcloudUserLikes struct {
	Collection []struct {
		CreatedAt time.Time `json:"created_at"`
		Track     struct {
			ID           uint   `json:"id"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			PermalinkURL string `json:"permalink_url"`
			ArtworkURL   string `json:"artwork_url"`
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
			id:          like.Track.ID,
			title:       like.Track.Title,
			description: like.Track.Description,
			pubdate:     like.CreatedAt,
			link:        like.Track.PermalinkURL,
			image:       like.Track.ArtworkURL,
			stream:      fmt.Sprintf("%s/%d.mp3", s.MediaSource, like.Track.ID),
		}
		tracks = append(tracks, &track)
	}

	return tracks, nil
}

type soundcloudStreams struct {
	URL string `json:"http_mp3_128_url"`
}

// StreamURL takes a Soundcloud song id and redirects the client to the real location of a song
func (s *SoundcloudSource) StreamURL(streamID string) string {
	streamsURL := fmt.Sprintf("https://api.soundcloud.com/tracks/%s/streams?client_id=%s", streamID, s.ClientID)
	resp, _ := s.Client.Get(streamsURL)
	streams := soundcloudStreams{}
	json.NewDecoder(resp.Body).Decode(&streams)
	return streams.URL
}
