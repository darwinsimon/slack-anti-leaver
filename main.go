package main

import (
	"log"
	"net/http"

	"github.com/nlopes/slack"
	c "github.com/robfig/cron"
)

func main() {

	// User token
	api := slack.New("SLACK TOKEN HERE")

	cronObj := c.New()

	// Channel ID
	channelID := "CHANNEL ID HERE"

	// Save timestamp
	lastTimestamp := ""

	// Hit API every 3 seconds
	cronObj.AddFunc("@every 3s", func() {

		historyParam := slack.HistoryParameters{
			Count:  30, // Message count per API call
			Oldest: lastTimestamp,
		}

		histories, err := api.GetChannelHistory(channelID, historyParam)
		if err != nil {
			log.Println("ERROR:", err)

			// Try to unarchive
			if err.Error() == "is_archived" {
				log.Println("Unarchive channel...")
				api.UnarchiveChannel(channelID)
			}

		} else {
			for _, message := range histories.Messages {
				if message.Msg.SubType == "channel_leave" {

					log.Println("Inviting", message.Msg.User, "to channel...")

					// Update timestamp
					if lastTimestamp < message.Msg.Timestamp {
						lastTimestamp = message.Msg.Timestamp
						_, err := api.InviteUserToChannel(channelID, message.Msg.User)
						if err != nil {
							log.Println("ERROR:", err)

							// Try to unarchive
							if err.Error() == "is_archived" {
								log.Println("Unarchive channel...")
								api.UnarchiveChannel(channelID)
							}
							break
						}

						log.Println("New timestamp", message.Msg.Timestamp)
					}
				}

			}
		}
	})

	cronObj.Start()

	log.Println("Started...")
	http.ListenAndServe(":3000", nil)
}
