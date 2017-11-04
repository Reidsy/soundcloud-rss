FROM golang:1.9 as build
WORKDIR /go/src/github.com/reidsy/soundcloud-rss/
RUN go get -d github.com/eduncan911/podcast
RUN go get -d gopkg.in/h2non/gock.v1
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a .

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/reidsy/soundcloud-rss/soundcloud-rss .
CMD ["./soundcloud-rss"]
EXPOSE 8080
