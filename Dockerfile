# syntax=docker/dockerfile:1.4

FROM golang:alpine as build
WORKDIR /app
COPY *.go /app/
RUN go mod init mastodon-group-bot && \
  go mod tidy && \
  go build -o mastodon-group-bot

FROM alpine:latest
WORKDIR /
COPY --chmod=555 entrypoint.sh /
COPY --from=build /app/mastodon-group-bot /

VOLUME ["/data"]
ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-config", "/data/config.json", "-db", "/data/limits.db" ]
