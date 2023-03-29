# syntax=docker/dockerfile:1
FROM golang:alpine
WORKDIR /app
COPY *.go ./
COPY config.json ./config.example.json
RUN go mod init mastodon-group-bot
RUN go mod tidy
RUN go build -o /mastodon-group-bot
RUN rm -rf /app
COPY --chmod=+x entrypoint.sh /

VOLUME ["/data"]
ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-config", "/data/config.json", "-db", "/data/limits.db" ]
