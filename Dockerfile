# syntax=docker/dockerfile:1.4

FROM golang:bullseye as build
WORKDIR /app
COPY *.go /app/
RUN go mod init mastodon-group-bot && \
  go mod tidy && \
  go build -o ./mastodon-group-bot

FROM debian:bullseye-slim
COPY --from=build /app/mastodon-group-bot ./
VOLUME ["/data"]
ENTRYPOINT ["/mastodon-group-bot"]
CMD [ "-config", "/data/config.json", "-db", "/data/limits.db", "-log", "/data/mastodon-group-bot.log" ]
