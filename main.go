package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
)

func main() {
	// Parse args
	ConfPath := flag.String("config", "config.json", "Path to config")
	flag.Parse()

	// Reading config
	type Config struct {
		Server         string   `json:"Server"`
		ClientID       string   `json:"ClientID"`
		ClientSecret   string   `json:"ClientSecret"`
		AccessToken    string   `json:"AccessToken"`
		WelcomeMessage string   `json:"WelcomeMessage"`
		Admins         []string `json:"Admins"`
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

		// Read message
		if notif.Type == "mention" {
			for i := 0; i < len(followers); i++ { // Follow check
				if notif.Status.Account.Acct == string(followers[i].Acct) {
					if notif.Status.Visibility == "public" { // Reblog toot
						c.Reblog(ctx, notif.Status.ID)
					} else if notif.Status.Visibility == "direct" { // Admin commands
						for y := 0; y < len(Conf.Admins); y++ {
							if notif.Status.Account.Acct == Conf.Admins[y] {
								text := notif.Status.Content
								recmd := regexp.MustCompile(`<.*?> `)
								command := recmd.ReplaceAllString(text, "")
								args := strings.Split(command, " ")

								if len(args) == 2 {
									switch args[0] {
									case "unboost":
										c.Unreblog(ctx, mastodon.ID((args[1])))
									case "delete":
										c.DeleteStatus(ctx, mastodon.ID(args[1]))
									case "block":
										c.AccountBlock(ctx, mastodon.ID(args[1]))
									case "unblock":
										c.AccountUnblock(ctx, mastodon.ID(args[1]))
									}
								}
							} else {
								var message = fmt.Sprintf("@%s%s", notif.Account.Acct, ", you are not admin!")
								postToot(message, "direct")
							}
						}
					}
				} else {
					var message = fmt.Sprintf("@%s%s", notif.Account.Acct, ", you are not subscribed!")
					postToot(message, "direct")
				}
			}
		}
	}
}
