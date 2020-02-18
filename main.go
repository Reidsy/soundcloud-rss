package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/profiler"
)

func main() {
	// start profiler
	profile()

	router, err := router()
	if err != nil {
		log.Fatal(err)
	}

	serveErr := http.ListenAndServe(":8080", router)
	log.Fatal(serveErr)
}

func router() (*http.ServeMux, error) {
	soundcloudAPI := SoundcloudSource{
		ClientID:    os.Getenv("SOUNDCLOUD_CLIENT_ID"),
		MediaSource: os.Getenv("MEDIA_SOURCE")}
	encoder := PodcastPlaylistEncoder{}
	feed := FeedServer{&soundcloudAPI, &encoder}
	stream := StreamServer{&soundcloudAPI}

	router := http.NewServeMux()
	router.HandleFunc("/monitoring/healthcheck", healthcheckEndpoint)
	router.Handle("/stream/", &stream)
	router.Handle("/", &feed)

	return router, nil
}

func healthcheckEndpoint(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s\n", http.StatusText(status))
}

func profile() {
	err := profiler.Start(profiler.Config{
		Service:        "com-reidsy-soundcloudrss",
		ProjectID:      "com-reidsy-soundcloudrss",
		DebugLogging:   false,
		MutexProfiling: true,
	})
	if err != nil {
		log.Printf("Failed to start Stackdriver Profiler: %v", err)
	}
}
