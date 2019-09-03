package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	gock "gopkg.in/h2non/gock.v1"
)

func TestSoundcloudSourceUser(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.soundcloud.com").
		Get("/resolve").
		MatchParam("url", "^https://soundcloud.com/reidsy$").
		Reply(200).
		JSON(soundcloudUser{
			ID:                   69198537,
			Permalink:            "reidsy",
			Username:             "Reidsy",
			URI:                  "https://api.soundcloud.com/users/69198537",
			PublicFavoritesCount: 102,
		})

	source := SoundcloudSource{
		ClientID: "not_my_client_id"}
	user, userErr := source.fetchUser("reidsy")

	if userErr != nil {
		t.Fatalf("userErr is not nil. Got %s", userErr)
	}
	if user.ID != 69198537 {
		t.Fatalf("User ID is incorrect. Got: %d", user.ID)
	}
	if user.Permalink != "reidsy" {
		t.Fatalf("Permalink is incorrect. Got: %s", user.Permalink)
	}
	if user.Username != "Reidsy" {
		t.Fatalf("Username is incorrect. Got: %s", user.Username)
	}
	if user.URI != "https://api.soundcloud.com/users/69198537" {
		t.Fatalf("URI is incorrect. Got: %s", user.URI)
	}
	if user.PublicFavoritesCount != 102 {
		t.Fatalf("PublicFavoritesCount is incorrect. Got: %d", user.PublicFavoritesCount)
	}

	if gock.IsDone() != true {
		t.Fatalf("Not all expected requests were made")
	}
}

func TestSoundcloudSourcePlaylist(t *testing.T) {
	user := soundcloudUser{
		ID:                   69198537,
		Permalink:            "reidsy",
		Username:             "Reidsy",
		URI:                  "https://api.soundcloud.com/users/69198537",
		PublicFavoritesCount: 102,
	}

	likes_json := `{
		"collection": [
			{
				"created_at": "2018-04-30T20:10:11Z",
				"track": {
					"id": 1,
					"title": "my song",
					"description": "it's really good",
					"permalink_url": "http://example.com/my-song",
					"artwork_url": "http://example.com/my-song.png"
				}
			},
			{
				"created_at": "2018-04-22T10:20:33Z",
				"track": {
					"id": 1,
					"title": "my song",
					"description": "",
					"permalink_url": "http://example.com/my-song",
					"artwork_url": "http://example.com/my-song.png"
				}
			}
		],
		"next_href": ""
	}`
	likes := soundcloudUserLikes{}
	json.NewDecoder(strings.NewReader(likes_json)).Decode(&likes)

	defer gock.Off()

	// stub user resolution
	gock.New("https://api.soundcloud.com").
		Get("/resolve").
		MatchParam("url", "^https://soundcloud.com/reidsy$").
		Reply(200).
		JSON(user)

	// stub likes endpoint
	gock.New("https://api-v2.soundcloud.com").
		Get("/users/69198537/likes").
		Reply(200).
		JSON(likes)

	source := SoundcloudSource{
		MediaSource: "http://example.com/streams",
		ClientID:    "my_client_id"}

	playlist, playlistErr := source.Playlist("reidsy", "likes")
	if playlistErr != nil {
		t.Fatalf("playlistErr is not nil. Got %s", playlistErr)
	}
	if playlist.Title() != "Likes" {
		t.Fatalf("Incorrect playlist title. Got: %s", playlist.Title())
	}
	if playlist.Link() != "https://soundcloud.com/reidsy/likes" {
		t.Fatalf("Incorrect playlist link. Got: %s", playlist.Link())
	}
	if playlist.Description() != "Likes" {
		t.Fatalf("Incorrect playlist description. Got: %s", playlist.Description())
	}
	if playlist.Author() != "Reidsy" {
		t.Fatalf("Incorrect playlist author. Got %s", playlist.Author())
	}
	if playlist.PubDate() == nil {
		t.Fatal("Playlist PubDate is not set")
	}
	if playlist.LastBuild() == nil {
		t.Fatal("Playlist LastBuild is not set")
	}
	if len(playlist.Tracks()) != 2 {
		t.Fatalf("Incorrect number of tracks in playlist. Got %d", len(playlist.Tracks()))
	}

	track := playlist.Tracks()[0]
	if track.ID() != 1 {
		t.Fatalf("Incorrect track ID. Got: %d", track.ID())
	}
	if track.Title() != "my song" {
		t.Fatalf("Incorrect track title. Got: %s", track.Title())
	}
	if track.Description() != "it's really good" {
		t.Fatalf("Incorrect track description. Got: %s", track.Title())
	}
	if track.PubDate().Format(time.RFC3339) != "2018-04-30T20:10:11Z" {
		t.Fatalf("Incorrect PubDate. Got: %s", track.PubDate())
	}
	if track.Link() != "http://example.com/my-song" {
		t.Fatalf("Incorrect track link. Got: %s", track.Link())
	}
	if track.Image() != "http://example.com/my-song.png" {
		t.Fatalf("Incorrect track image. Got: %s", track.Image())
	}
	if track.Stream() != "http://example.com/streams/1.mp3" {
		t.Fatalf("Incorrect track stream. Got: %s", track.Stream())
	}

	trackWithoutDescription := playlist.Tracks()[1]
	if trackWithoutDescription.Description() != trackWithoutDescription.Title() {
		t.Fatal("Title should be used in place of description when description is not set")
	}
}

func TestSoundcloudSourceStreamURL(t *testing.T) {
	defer gock.Off()

	// stub streams endpoint
	gock.New("https://api.soundcloud.com").
		Get("/tracks/1234/streams").
		Reply(200).
		JSON(map[string]string{"http_mp3_128_url": "https://example.com/mysong.mp3"})

	source := SoundcloudSource{
		ClientID: "my_client_id"}

	streamURL := source.StreamURL("1234")
	if streamURL != "https://example.com/mysong.mp3" {
		t.Fatalf("Unexpected streamURL. Got: %s", streamURL)
	}

	if gock.IsDone() != true {
		t.Fatalf("Not all expected requests were made")
	}
}
