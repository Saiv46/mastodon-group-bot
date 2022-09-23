package main

import "github.com/mattn/go-mastodon"

// Posting function
func postToot(toot string, vis string) (*mastodon.Status, error) {
	conToot := mastodon.Toot{
		Status:     toot,
		Visibility: vis,
	}
	status, err := c.PostStatus(ctx, &conToot)
	return status, err
}
