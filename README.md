# Mastodon group bot
This is a bot which implements group functionality in Mastodon.

## Features
* Repost toots
* Welcome message of new members
* Limit of toots per hour
* Duplicate protection
* Order limit
* Notification cleaning
* Logging
* Admin commands

### Admin commands
* unboost \<Toot ID>
* delete  \<Toot ID>

# Configuration
The bot is configured in a JSON file that looks like this:
```
{
    "Server":               "https://example.com",
    "ClientID":             "0000000000000000000000000000000000000000000",
    "ClientSecret":         "0000000000000000000000000000000000000000000",
    "AccessToken":          "0000000000000000000000000000000000000000000",
    "WelcomeMessage":       "We have a new member in our group. Please love and favor",
    "NotFollowedMessage":   "you are not followed",
    "Max_toots":            2,
    "Toots_interval":       12,
    "Duplicate_buf":        10,
    "Order_limit":          1,
    "Del_notices_interval": 30,
    "Admins":               ["admin@example.com"]
}
```

# Building
```
go mod init mastodon-group-bot
go mod tidy
go build
```

# Setup services
For first make dirs, copy config and binary
```
mkdir /etc/mastodon-group-bot
mkdir /var/lib/mastodon-group-bot
mkdir /var/log/mastodon-group-bot
chown nobody /var/lib/mastodon-group-bot
chown nobody /var/log/mastodon-group-bot
cp config.json /etc/mastodon-group-bot/config.json
cp mastodon-group-bot /usr/bin/mastodon-group-bot
```

## Systemd
```
cp ./services/systemd/mastodon-group-bot.service /etc/systemd/system/mastodon-group-bot.service
```

## OpenRC
```
cp ./services/openrc/mastodon-group-bot /etc/init.d/mastodon-group-bot
```

# Usage
```
mastodon-group-bot -config <path> -db <path> -log <path>
```