package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"gitlab.xmltravelgate.com/core/core-event-router/models"
	"gitlab.xmltravelgate.com/core/core-event-router/package/pubsubService"
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

	ctx := context.Background()
	pbService := pubsubService.NewPubSubService(ctx, PROJECT_ID, TOPIC_NAME, PRODUCER_CREDENTIALS, EVENT_ROUTER_CLIENTID)

	//create example message
	msg := createExampleMsg()
	pbService.NotifyEvent(&msg)

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

func createExampleMsg() models.MessageModel {
	attr := models.MessageAttributes{
		Src:   "IAM",
		Prod:  "fake_prod",
		Type:  "fake_type",
		Stype: "fake_stype",
		Op:    "fake_op",
	}
	return models.MessageModel{
		Payload:         "fake_payload",
		SpecificPayload: "fake_specific_payload",
		Attributes:      attr,
	}
}
