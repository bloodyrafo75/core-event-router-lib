package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bloodyrafo75/core-event-router-lib/models"
	"github.com/bloodyrafo75/core-event-router-lib/package/pubsubService"
	"github.com/joho/godotenv"
)

var (
	PUBSUB_HOST           string //Only for working with local docker
	PROJECT_ID            string
	TOPIC_NAME            string
	PORT                  string //Only for working with local docker
	EVENT_ROUTER_CLIENTID string
	SUBSCRIPTION_ID       string
	LOCAL_MODE            string
	CONSUMER_CREDENTIALS  string
)

func main() {
	err := getEnvConfiguration()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pbService := pubsubService.NewPubSubService(ctx, PROJECT_ID, TOPIC_NAME, CONSUMER_CREDENTIALS, EVENT_ROUTER_CLIENTID)
	callBackFn := callback()

	err = pbService.StartConsumer((*pubsubService.Callback)(&callBackFn), SUBSCRIPTION_ID)
	if err != nil {
		fmt.Println(err)
	}
}

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
	SUBSCRIPTION_ID = os.Getenv("SUBSCRIPTION_ID")
	CONSUMER_CREDENTIALS = os.Getenv("CONSUMER_CREDENTIALS")
	LOCAL_MODE = os.Getenv("LOCAL_MODE")

	if LOCAL_MODE == "1" {
		err := os.Setenv("PUBSUB_EMULATOR_HOST", PUBSUB_HOST) // Set PUBSUB_EMULATOR_HOST environment variable.
		if err != nil {
			return err
		}
	}

	return nil
}

func callback() func(*models.MessageModel) error {
	return func(msg *models.MessageModel) error {
		if msg == (&models.MessageModel{}) {
			return fmt.Errorf("response doesn't fit to model %v", msg)
		}

		fmt.Println("Process msg:")
		fmt.Println(msg)
		return nil
	}
}
