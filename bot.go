package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
)

func run_bot(Conf Config, DB string) {
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

	my_account, err := c.GetAccountCurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	followers, err := c.GetAccountFollowers(ctx, my_account.ID, &mastodon.Pagination{Limit: 60})
	if err != nil {
		log.Fatal(err)
	}

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

		// New follower
		if notif.Type == "follow" {
			acct := notif.Account.Acct
			if !followed(acct, DB) { // Add to db and post welcome message
				add_to_db(acct, Conf.Max_toots, DB)
				var message = fmt.Sprintf("%s @%s", Conf.WelcomeMessage, acct)
				postToot(message, "public")
			}
		}

		// Read message
		if notif.Type == "mention" {
			acct := notif.Status.Account.Acct
			for i := 0; i < len(followers); i++ {
				if acct == string(followers[i].Acct) { // Follow check
					if notif.Status.Visibility == "public" { // Reblog toot
						if notif.Status.InReplyToID == nil { // Not boost replies
							if !followed(acct, DB) { // Add to db if needed
								add_to_db(acct, Conf.Max_toots, DB)
							}
							if check_ticket(acct, Conf.Max_toots, Conf.Toots_interval, DB) > 0 { // Limit
								take_ticket(acct, DB)
								c.Reblog(ctx, notif.Status.ID)
							}
						}
					} else if notif.Status.Visibility == "direct" { // Admin commands
						for y := 0; y < len(Conf.Admins); y++ {
							if acct == Conf.Admins[y] {
								text := notif.Status.Content
								recmd := regexp.MustCompile(`<.*?> `)
								command := recmd.ReplaceAllString(text, "")
								args := strings.Split(command, " ")
								mID := mastodon.ID((args[1]))

								if len(args) == 2 {
									switch args[0] {
									case "unboost":
										c.Unreblog(ctx, mID)
									case "delete":
										c.DeleteStatus(ctx, mID)
									case "block":
										c.AccountBlock(ctx, mID)
									case "unblock":
										c.AccountUnblock(ctx, mID)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
