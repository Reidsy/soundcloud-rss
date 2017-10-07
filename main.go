package main

import (
	"fmt"
	"net/http"
)

func main() {
	srss := SoundcloudRSSServer{}
	http.HandleFunc("/monitoring/healthcheck", healthcheck)
	http.Handle("/", &srss)
	http.ListenAndServe(":8080", nil)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s\n", http.StatusText(status))
}
