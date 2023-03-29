#!/bin/bash

# The application log will be redirected to the main docker container process's stdout, so
# that it will show up in the container logs
touch /var/log/mastodon-group-bot.log
ln -sf /proc/1/fd/1 /var/log/mastodon-group-bot.log

/mastodon-group-bot -log /var/log/mastodon-group-bot.log "$@"
