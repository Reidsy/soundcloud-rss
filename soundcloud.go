package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SoundcloudAPI struct {
	Client       http.Client
	ClientID     string
	ClientSecret string
}

type SoundcloudAPIPlaylist struct {
	Title       string
	Link        string
	Description string
	Author      string
	Tracks      []SoundcloudAPITrack
}

func (s *SoundcloudAPI) getFavouritesPlaylist(username string) SoundcloudAPIPlaylist {
	user := s.getUser(username)
	title := fmt.Sprintf("%s's Likes", user.Username)
	likesURL := fmt.Sprintf("https://soundcloud.com/%s/likes", user.Permalink)
	likes := s.getLikes(user)
	tracks := []SoundcloudAPITrack{}
	for _, like := range likes {
		tracks = append(tracks, like.Track)
	}
	playlist := SoundcloudAPIPlaylist{
		Title:       title,
		Link:        likesURL,
		Description: title,
		Author:      strings.Title(username),
		Tracks:      tracks,
	}
	return playlist
}

func (s *SoundcloudAPI) getPlaylist(username string, playlist string) {

}

func (s *SoundcloudAPI) resolveUrl(asd string) string {
	u, _ := url.Parse("https://api.soundcloud.com/resolve")
	q := u.Query()
	q.Set("client_id", s.ClientID)
	q.Set("url", asd)
	u.RawQuery = q.Encode()
	return u.String()
}

// SoundcloudAPIUser @TODO
type SoundcloudAPIUser struct {
	ID                   uint   `json:"id"`
	Permalink            string `json:"permalink"`
	Username             string `json:"username"`
	URI                  string `json:"uri"`
	PublicFavoritesCount uint   `json:"public_favorites_count"`
}

func (s *SoundcloudAPI) getUser(username string) SoundcloudAPIUser {
	u, _ := url.Parse("https://api.soundcloud.com/resolve")
	q := u.Query()
	q.Set("client_id", s.ClientID)
	q.Set("url", fmt.Sprintf("https://soundcloud.com/%s", username))
	u.RawQuery = q.Encode()

	resp, _ := s.Client.Get(u.String())

	user := SoundcloudAPIUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	return user
}

// SoundcloudAPITrack @TODO
type SoundcloudAPITrack struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"permalink_url"`
}

// SoundcloudAPILike @TODO
type SoundcloudAPILike struct {
	CreatedAt time.Time          `json:"created_at"`
	Track     SoundcloudAPITrack `json:"track"`
}

type SoundcloudAPILikeRequest struct {
	Collection []SoundcloudAPILike `json:"collection"`
	NextHref   string              `json:"next_href"`
}

func (s *SoundcloudAPI) getLikes(user SoundcloudAPIUser) []SoundcloudAPILike {
	u, _ := url.Parse(fmt.Sprintf("https://api-v2.soundcloud.com/users/%d/likes", user.ID))

	q := u.Query()
	q.Set("offset", fmt.Sprint(0))
	q.Set("limit", fmt.Sprint(user.PublicFavoritesCount))
	q.Set("client_id", s.ClientID)
	u.RawQuery = q.Encode()

	resp, _ := s.Client.Get(u.String())

	likeRequest := SoundcloudAPILikeRequest{}
	json.NewDecoder(resp.Body).Decode(&likeRequest)
	return likeRequest.Collection
}

// SoundcloudAPIStreams @TODO
type SoundcloudAPIStreams struct {
	URL string `json:"http_mp3_128_url"`
}

func (s *SoundcloudAPI) getStreamURL(streamID string) string {
	streamsURL := fmt.Sprintf("https://api.soundcloud.com/tracks/%s/streams?client_id=%s", streamID, s.ClientID)
	resp, _ := s.Client.Get(streamsURL)
	streams := SoundcloudAPIStreams{}
	json.NewDecoder(resp.Body).Decode(&streams)
	return streams.URL
}
