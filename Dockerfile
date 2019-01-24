FROM golang:1.11 as build
WORKDIR /go/src/github.com/reidsy/soundcloud-rss/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/reidsy/soundcloud-rss/soundcloud-rss /go/bin/soundcloud-rss
CMD ["/go/bin/soundcloud-rss"]
EXPOSE 8080
