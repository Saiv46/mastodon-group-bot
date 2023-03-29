# syntax=docker/dockerfile:1
FROM golang:alpine
WORKDIR /app
COPY *.go ./
COPY config.json ./config.example.json
RUN go mod init mastodon-group-bot
RUN go mod tidy
RUN go build -o /mastodon-group-bot

WORKDIR /
COPY --chmod=+x entrypoint.sh ./
RUN rm -rf /app

VOLUME ["/data"]
ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-config", "/data/config.json", "-db", "/data/limits.db" ]
