package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	soundcloudAPI := SoundcloudSource{
		ClientID:    os.Getenv("SOUNDCLOUD_CLIENT_ID"),
		MediaSource: os.Getenv("MEDIA_SOURCE")}
	encoder := PodcastPlaylistEncoder{}
	feed := FeedServer{&soundcloudAPI, &encoder}
	stream := StreamServer{&soundcloudAPI}

	http.HandleFunc("/monitoring/healthcheck", healthcheck)
	http.Handle("/stream/", &stream)
	http.Handle("/", &feed)
	http.ListenAndServe(":8080", nil)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s\n", http.StatusText(status))
}
