package main

import (
	"fmt"
	"log"
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

	router := http.NewServeMux()
	router.HandleFunc("/monitoring/healthcheck", healthcheck)
	router.Handle("/stream/", &stream)
	router.Handle("/", &feed)

	serveErr := http.ListenAndServe(":8080", router)
	log.Fatal(serveErr)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s\n", http.StatusText(status))
}
