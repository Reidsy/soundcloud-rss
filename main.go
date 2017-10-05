package main

import "net/http"

func main() {
	srss := SoundcloudRSSServer{}
	http.Handle("/", &srss)
	http.ListenAndServe(":8080", nil)
}
