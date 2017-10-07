package main

// SoundcloudRSSServer @TODO
import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/eduncan911/podcast"
)

// SoundcloudRSSServer @TODO
type SoundcloudRSSServer struct {
}

func (srss *SoundcloudRSSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		srss.sendStatus(w, http.StatusMethodNotAllowed)
		return
	}
	srssr := SoundcloudRSSRequest{}
	if srssr.ParseURL(r.URL) == ErrMalformedURL {
		fmt.Printf("%v\n", r.URL)
		srss.sendStatus(w, http.StatusNotFound)
		return
	}
	if srssr.IsFeed() {
		srss.feed(w, srssr.Username, srssr.Playlist)
	} else if srssr.IsStream() {
		srss.stream(w, srssr.StreamID())
	}
}

func (srss *SoundcloudRSSServer) sendStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s\n", http.StatusText(status))
}

func (srss *SoundcloudRSSServer) getClientId() string {
	return os.Getenv("SOUNDCLOUD_CLIENT_ID")
}

func (srss *SoundcloudRSSServer) feed(w http.ResponseWriter, username string, playlist string) {
	s := SoundcloudAPI{ClientID: srss.getClientId()}
	pubDate := time.Now()
	plist := s.getFavouritesPlaylist(username)
	pcast := podcast.New(plist.Title, plist.Link, plist.Description, &pubDate, &pubDate)
	for idx, track := range plist.Tracks {
		itemPubDate := pubDate.AddDate(0, 0, -idx)
		pitem := podcast.Item{
			Title:       track.Title,
			Description: track.Description,
			PubDate:     &itemPubDate,
			Link:        track.Link,
		}
		if pitem.Description == "" {
			pitem.Description = pitem.Title
		}
		streamURL := fmt.Sprintf("http://soundcloudrss.reidsy.com/%s/%s/%d.mp3", username, playlist, track.ID)
		pitem.AddEnclosure(streamURL, podcast.MP3, 0)
		_, err := pcast.AddItem(pitem)
		if err != nil {
			fmt.Printf("addItemErr: (%d) %s", idx, err)
		}
	}
	pcast.Encode(w)
}

func (srss *SoundcloudRSSServer) stream(w http.ResponseWriter, streamID string) {
	s := SoundcloudAPI{ClientID: srss.getClientId()}
	url := s.getStreamURL(streamID)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}
