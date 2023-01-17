package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bloodyrafo75/core-event-router-lib/package/pubsubService"
	"github.com/joho/godotenv"
)

var (
	PUBSUB_HOST           string //Only for working with local docker
	PROJECT_ID            string
	TOPIC_NAME            string
	PORT                  string //Only for working with local docker
	EVENT_ROUTER_CLIENTID string
	LOCAL_MODE            string
	PRODUCER_CREDENTIALS  string
)

func main() {
	err := getEnvConfiguration()
	if err != nil {
		panic(err)
	}

	pubsubService.NewPubSubService(context.Background(), PROJECT_ID, TOPIC_NAME, PRODUCER_CREDENTIALS, EVENT_ROUTER_CLIENTID)

	//create example message
	msg := pubsubService.PubsubClient.CreatePubsubMsg("IAM", "organizations", "READ", "ORG", "update", "", "")
	res, err := pubsubService.PubsubClient.NotifyEvent(&msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}

// get configuration from .env file.
func getEnvConfiguration() error {
	err := godotenv.Load("configs/.env")
	if err != nil {
		return err
	}

	PUBSUB_HOST = os.Getenv("PUBSUB_HOST")
	PROJECT_ID = os.Getenv("PROJECT_ID")
	TOPIC_NAME = os.Getenv("TOPIC_NAME")
	PORT = os.Getenv("PORT")
	EVENT_ROUTER_CLIENTID = os.Getenv("EVENT_ROUTER_CLIENTID")
	PRODUCER_CREDENTIALS = os.Getenv("PRODUCER_CREDENTIALS")
	LOCAL_MODE = os.Getenv("LOCAL_MODE")

	if LOCAL_MODE == "1" {
		err := os.Setenv("PUBSUB_EMULATOR_HOST", PUBSUB_HOST) // Set PUBSUB_EMULATOR_HOST environment variable.
		if err != nil {
			return err
		}
	}

	return nil
}
