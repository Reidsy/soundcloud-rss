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

- `SOUNDCLOUD_CLIENT_ID` can be obtained by inspecting an api request from the soundcloud web app and looking for the `client_id=` query param.
- `MEDIA_SOURCE` is where the server will be running, by default this will be localhost on port 8080.

### Running

```
$ docker-compose up
```

## Testing and Development

### Requirements

- [go 1.13](https://golang.org/dl)

### Development

Run the following commands before opening a pull request.

```
$ go fmt
$ go vet
```

### Running Tests

```
$ go test
```

### Compiling

```
$ go build
```
