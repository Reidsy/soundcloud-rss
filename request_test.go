package main

import (
	"net/url"
	"testing"
)

func TestRequestParseURLFeed(t *testing.T) {
	rssURL, _ := url.Parse("http://reidsy.com/reidsy/likes/rss.xml")
	r := SoundcloudRSSRequest{}
	err := r.ParseURL(rssURL)
	if err != nil {
		t.Fatal("Failed to parse good url")
	}
	if r.Username != "reidsy" {
		t.Fatal("Username != reidsy")
	}
	if r.Playlist != "likes" {
		t.Fatal("Playlist != likes")
	}
	if r.Item != "rss.xml" {
		t.Fatal("Item != rss.xml")
	}
}

func TestRequestParseURLStream(t *testing.T) {
	rssURL, _ := url.Parse("http://reidsy.com/reidsy/likes/11254.mp3")
	r := SoundcloudRSSRequest{}
	err := r.ParseURL(rssURL)
	if err != nil {
		t.Fatal("Failed to parse good url")
	}
	if r.Username != "reidsy" {
		t.Fatal("Username != reidsy")
	}
	if r.Playlist != "likes" {
		t.Fatal("Playlist != likes")
	}
	if r.Item != "11254.mp3" {
		t.Fatal("Item != 11254.mp3")
	}
}

func TestRequestParseURLBad(t *testing.T) {
	rssURL, _ := url.Parse("http://reidsy.com/reidsy/likes/hello.world")
	r := SoundcloudRSSRequest{}
	err := r.ParseURL(rssURL)
	if err != ErrMalformedURL {
		t.Fatal("Failed to detect bad url")
	}
	if r.Username != "" {
		t.Fatal("Username has been set")
	}
	if r.Playlist != "" {
		t.Fatal("Playlist has been set")
	}
	if r.Item != "" {
		t.Fatal("Item has been set")
	}
}

func TestRequestIsFeed(t *testing.T) {
	isFeed := SoundcloudRSSRequest{Item: "rss.xml"}
	if isFeed.IsFeed() != true {
		t.Fatalf("isFeed incorrectly identified %s", isFeed.Item)
	}

	notFeed := SoundcloudRSSRequest{Item: "123.mp3"}
	if notFeed.IsFeed() != false {
		t.Fatalf("notFeed incorrectly identified %s", isFeed.Item)
	}

	badFeed := SoundcloudRSSRequest{Item: "fake.item"}
	if badFeed.IsFeed() != false {
		t.Fatalf("badFeed incorrectly identified %s", isFeed.Item)
	}
}

func TestRequestIsStream(t *testing.T) {
	isStream := SoundcloudRSSRequest{Item: "123.mp3"}
	if isStream.IsStream() != true {
		t.Fatalf("isStream incorrectly identified %s", isStream.Item)
	}

	notStream := SoundcloudRSSRequest{Item: "rss.xml"}
	if notStream.IsStream() != false {
		t.Fatalf("notStream incorrectly identified %s", notStream.Item)
	}

	badStream := SoundcloudRSSRequest{Item: "fake.item"}
	if badStream.IsStream() != false {
		t.Fatalf("badStream incorrectly identified %s", badStream.Item)
	}
}

func TestRequestStreamId(t *testing.T) {
	isStream := SoundcloudRSSRequest{Item: "123.mp3"}
	if isStream.StreamID() != "123" {
		t.Fatalf("isStream did not find streamID %s", isStream.Item)
	}

	notStream := SoundcloudRSSRequest{Item: "rss.xml"}
	if notStream.StreamID() != "" {
		t.Fatalf("notStream found incorrect streamID %s", notStream.Item)
	}
}
