# Mastodon group bot

This is a bot which implements group functionality in Mastodon.

# Configuration

The bot is configured in a JSON file that looks like this:

```
{
    "Server":           "https://example.com",
    "ClientID":         "0000000000000000000000000000000000000000000",
    "ClientSecret":     "0000000000000000000000000000000000000000000",
    "AccessToken":      "0000000000000000000000000000000000000000000",
    "WelcomeMessage":   "We have a new member in our group. Please love and favor"
    "Max_toots":        1,
    "Toots_interval":   24,
    "Admins":           ["admin@example.com"]
}
```

# Building

```
go mod init mastodon-group-bot

go mod tidy

go build
```

# Usage

```
Usage of mastodon-group-bot:
  -config string
        Path to config (default "config.json")
```