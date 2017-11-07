# Souncloud-RSS

Create a podcast feed from your Soundcloud playlists and likes.

Feeds are served in the format:

```
http://localhost:8080/:username/:playlist/rss.xml
```

Examples:

- http://localhost:8080/reidsy/likes/rss.xml
- http://localhost:8080/reidsy/umf-2017/rss.xml


## Quickstart

### Requirements
- [Docker](http://docker.com)

### Configuration

Create a configuration file `config.env` with the following keys:
```
SOUNDCLOUD_CLIENT_ID=${CLIENT_ID}
MEDIA_SOURCE=http://localhost:8080/stream
```

- `SOUNDCLOUD_DOT_COM_CLIENT_ID` can be obtained by inspecting an api request from the soundcloud web app and looking for the `client_id=` query param.
- `MEDIA_SOURCE` is where the server will be running, by default this will be localhost on port 8080.

### Running

```
$ docker-compose up
```

## Testing and Development

### Requirements

- [go 1.9](https://golang.org/dl)
- [github.com/eduncan911/podcast](https://github.com/eduncan911/podcast)
- [gopkg.in/h2non/gock.v1](https://github.com/h2non/gock/tree/v1.0.6)

### Dependency Installation

```
$ go get -d github.com/eduncan911/podcast
$ go get -d gopkg.in/h2non/gock.v1
```

### Running Tests

```
$ go test
```

### Compiling

```
$ go build
```