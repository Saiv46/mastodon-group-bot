package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
)

func RunBot(Conf Config) {
	logger_init()

	c := mastodon.NewClient(&mastodon.Config{
		Server:       Conf.Server,
		ClientID:     Conf.ClientID,
		ClientSecret: Conf.ClientSecret,
		AccessToken:  Conf.AccessToken,
	})

	ctx := context.Background()
	events, err := c.StreamingUser(ctx)
	if err != nil {
		ErrorLogger.Println("Streaming")
	}

	my_account, err := c.GetAccountCurrentUser(ctx)
	if err != nil {
		ErrorLogger.Println("Fetch account info")
	}
	followers, err := c.GetAccountFollowers(ctx, my_account.ID, &mastodon.Pagination{Limit: 60})
	if err != nil {
		ErrorLogger.Println("Fetch followers")
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
			if !followed(acct) { // Add to db and post welcome message
				InfoLogger.Printf("%s followed", acct)

				add_to_db(acct, Conf.Max_toots)
				InfoLogger.Printf("%s added to database", acct)

				var message = fmt.Sprintf("%s @%s", Conf.WelcomeMessage, acct)
				err := postToot(message, "public")
				if err != nil {
					ErrorLogger.Println("Post welcome message")
				}
				InfoLogger.Printf("%s was welcomed", acct)
			}
		}

		// Read message
		if notif.Type == "mention" {
			acct := notif.Status.Account.Acct
			for i := 0; i < len(followers); i++ {
				if acct == string(followers[i].Acct) { // Follow check
					if notif.Status.Visibility == "public" { // Reblog toot
						if notif.Status.InReplyToID == nil { // Not boost replies
							if !followed(acct) { // Add to db if needed
								add_to_db(acct, Conf.Max_toots)
								InfoLogger.Printf("%s added to database", acct)
							}
							if check_ticket(acct, Conf.Max_toots, Conf.Toots_interval) > 0 { // Limit
								take_ticket(acct)
								InfoLogger.Printf("Ticket of %s was taken", acct)
								c.Reblog(ctx, notif.Status.ID)
								InfoLogger.Printf("Toot %s of %s was rebloged", notif.Status.URL, acct)
							} else {
								WarnLogger.Printf("%s haven't tickets", acct)
							}
						} else {
							WarnLogger.Printf("%s is reply and not boosted", notif.Status.URL)
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
							} else {
								break
							}
						}
					} else {
						WarnLogger.Printf("%s is not public toot and not boosted", notif.Status.URL)
						break
					}
				}
			}
		}
	}
}
