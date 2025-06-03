package main

import "log"

func sendNotification(msg string) {
	switch *notify {
	case "twitter":
		if _, err := tweet(msg); err != nil {
			log.Print(err)
		}
	case "slack":
		if err := slackNotify(msg); err != nil {
			log.Print(err)
		}
	default: // both
		if _, err := tweet(msg); err != nil {
			log.Print(err)
		}
		if err := slackNotify(msg); err != nil {
			log.Print(err)
		}
	}
}
