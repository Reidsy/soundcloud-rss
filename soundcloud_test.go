package main

import "testing"

func TestSoundcloudAPIResolveUrl(t *testing.T) {
	expectedURL := "https://api.soundcloud.com/resolve?client_id=abc123&url=this_is_a_test"
	s := SoundcloudAPI{ClientID: "abc123"}
	resolvedURL := s.resolveUrl("this_is_a_test")
	if resolvedURL != expectedURL {
		t.Fatalf("resolvedURL is incorrect. Expected: %s, Got: %s", expectedURL, resolvedURL)
	}
}
