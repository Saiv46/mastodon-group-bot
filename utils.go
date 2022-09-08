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

// Check following
func check_following(followers []*mastodon.Account, acct string) bool {
	for i := range followers {
		if acct == string(followers[i].Acct) {
			return true
		}
	}
	return false
}
