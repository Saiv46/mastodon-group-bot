package main

import (
	"sync"
	"time"

	"github.com/mattn/go-mastodon"
)

var (
	wg sync.WaitGroup
)

// Delete notices
func DeleteNotices() {
	wg.Done()

	for {
		statuses, err := c.GetAccountStatuses(ctx, my_account.ID, &mastodon.Pagination{Limit: 60})
		if err != nil {
			ErrorLogger.Println("Get account statuses")
		}

		time.Sleep(time.Duration(Conf.Del_notices_interval) * time.Second)

		for i := range statuses {
			if statuses[i].Visibility == "direct" {
				c.DeleteStatus(ctx, statuses[i].ID)
			}
		}
		InfoLogger.Println("Cleaning notices")

		reset_notice_counter()
		InfoLogger.Println("Reset notice counter")
	}
}
