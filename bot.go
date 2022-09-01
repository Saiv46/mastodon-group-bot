package main

import (
	"context"
	"crypto/sha512"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
)

func RunBot() {
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

				add_to_db(acct)
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
			content := notif.Status.Content
			tooturl := notif.Status.URL

			for i := 0; i < len(followers); i++ {
				if acct == string(followers[i].Acct) { // Follow check
					if notif.Status.Visibility == "public" { // Reblog toot
						if notif.Status.InReplyToID == nil { // Not boost replies
							// Duplicate protection
							content_hash := sha512.New()
							content_hash.Write([]byte(content))
							hash := fmt.Sprintf("%x", content_hash.Sum(nil))

							if !check_msg_hash(hash) {
								save_msg_hash(hash)
								InfoLogger.Printf("Hash of %s added to database", tooturl)
							} else {
								WarnLogger.Printf("%s is a duplicate and not boosted", tooturl)
								break
							}

							// Add to db if needed
							if !followed(acct) {
								add_to_db(acct)
								InfoLogger.Printf("%s added to database", acct)
							}

							// Message limit
							if check_ticket(acct) > 0 {
								take_ticket(acct)
								InfoLogger.Printf("Ticket of %s was taken", acct)
								c.Reblog(ctx, notif.Status.ID)
								InfoLogger.Printf("Toot %s of %s was rebloged", tooturl, acct)
							} else {
								WarnLogger.Printf("%s haven't tickets", acct)
							}
						} else {
							WarnLogger.Printf("%s is reply and not boosted", tooturl)
						}
					} else if notif.Status.Visibility == "direct" { // Admin commands
						for y := 0; y < len(Conf.Admins); y++ {
							if acct == Conf.Admins[y] {
								recmd := regexp.MustCompile(`<.*?> `)
								command := recmd.ReplaceAllString(content, "")
								args := strings.Split(command, " ")
								mID := mastodon.ID((args[1]))

								if len(args) == 2 {
									switch args[0] {
									case "unboost":
										c.Unreblog(ctx, mID)
									case "delete":
										c.DeleteStatus(ctx, mID)
									}
								}
							} else {
								break
							}
						}
					} else {
						WarnLogger.Printf("%s is not public toot and not boosted", tooturl)
						break
					}
				}
			}
		}
	}
}
