package main

import (
	"os"

	"github.com/timonv/pusher"
)

func setupClient() *pusher.Client {
	appID := os.Getenv("HUSTLE_TEST_APP_ID")
	key := os.Getenv("HUSTLE_TEST_KEY")
	secret := os.Getenv("HUSTLE_TEST_SECRET")

	client := pusher.NewClient(appID, key, secret)
	client.Host = os.Getenv("HUSTLE_HTTPPUBADDR")
	return client
}
