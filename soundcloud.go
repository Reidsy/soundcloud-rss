package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

type SoundcloudAPITrack struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"permalink_url"`
}

func (s *SoundcloudAPI) getFavouritesPlaylist(username string) SoundcloudAPIPlaylist {
	title := fmt.Sprintf("%s's Likes", strings.Title(username))
	likesURL := fmt.Sprintf("https://soundcloud.com/%s/likes", username)
	likesPlaylistURL := s.resolveUrl(likesURL)
	resp, _ := s.Client.Get(likesPlaylistURL)
	playlist := SoundcloudAPIPlaylist{
		Title:       title,
		Link:        likesURL,
		Description: title,
		Author:      strings.Title(username),
	}
	json.NewDecoder(resp.Body).Decode(&playlist.Tracks)
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

func (s *SoundcloudAPI) streamTrack(streamID string, w io.Writer) {
	streamURL := fmt.Sprintf("https://api.soundcloud.com/tracks/%s/stream?client_id=%s", streamID, s.ClientID)
	resp, _ := s.Client.Get(streamURL)
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
}
