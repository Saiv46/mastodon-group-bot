package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-mastodon"
)

func main() {
	// Parse args
	ConfPath := flag.String("config", "config.json", "Path to config")
	flag.Parse()

	// Reading config
	type Config struct {
		Server         string `json:"Server"`
		ClientID       string `json:"ClientID"`
		ClientSecret   string `json:"ClientSecret"`
		AccessToken    string `json:"AccessToken"`
		WelcomeMessage string `json:"WelcomeMessage"`
	}

	data, err := os.ReadFile(*ConfPath)
	if err != nil {
		log.Fatal(err)
	}

	var Conf Config
	json.Unmarshal(data, &Conf)

	c := mastodon.NewClient(&mastodon.Config{
		Server:       Conf.Server,
		ClientID:     Conf.ClientID,
		ClientSecret: Conf.ClientSecret,
		AccessToken:  Conf.AccessToken,
	})

	ctx := context.Background()
	events, err := c.StreamingUser(ctx)
	if err != nil {
		log.Fatal(err)
	}

	my_account, _ := c.GetAccountCurrentUser(ctx)
	followers, _ := c.GetAccountFollowers(ctx, my_account.ID, &mastodon.Pagination{Limit: 60})
	signed := false

	// Run bot
	for {
		notifEvent, ok := (<-events).(*mastodon.NotificationEvent)
		if !ok {
			continue
		}

		notif := notifEvent.Notification

		// Posting function
		postToot := func(toot string, vis string) error {
			conToot := mastodon.Toot{
				Status:     toot,
				Visibility: vis,
			}
			_, err := c.PostStatus(ctx, &conToot)
			return err
		}

		// New subscriber
		if notif.Type == "follow" {
			var message = fmt.Sprintf("%s @%s", Conf.WelcomeMessage, notif.Account.Acct)
			postToot(message, "public")
		}

		// Reblog toot
		if notif.Type == "mention" {
			sender := notif.Status.Account.Acct

			// Subscription check
			for i := 0; i < len(followers); i++ {
				if sender == string(followers[i].Acct) {
					signed = true
				}
			}

			if signed {
				c.Reblog(ctx, notif.Status.ID)
			} else {
				var message = fmt.Sprintf("@%s%s", notif.Account.Acct, ", you are not subscribed!")
				postToot(message, "direct")
			}
		}
	}
}
