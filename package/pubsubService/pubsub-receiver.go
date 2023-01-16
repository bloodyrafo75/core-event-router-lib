package pubsubService

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/mitchellh/mapstructure"
	"gitlab.xmltravelgate.com/core/core-event-router/models"
)

type PubSubReceiver struct {
	ctx            context.Context
	connection     PubSubConnection
	subscriptionId string
	callback       Callback
}

func NewPubSubReceiver(ctx context.Context, connection PubSubConnection, subscriptionId string, callback Callback) *PubSubReceiver {
	return &PubSubReceiver{
		ctx:            ctx,
		connection:     connection,
		subscriptionId: subscriptionId,
		callback:       callback,
	}
}

func (r *PubSubReceiver) Start() error {
	client, err := r.connection.Connect()
	if err != nil {
		return err
	}

	err = r.retrieveMessages(*client, r.subscriptionId)
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

func (r *PubSubReceiver) retrieveMessages(client pubsub.Client, subscriptionId string) error {
	sub := client.Subscription(subscriptionId)
	cctx, cancel := context.WithCancel(r.ctx)
	defer cancel()

	err := sub.Receive(cctx, func(_ context.Context, msg *pubsub.Message) {
		msg.Ack()
		msgModel := convertToTGXModel(msg)
		r.callback(msgModel)
	})
	if err != nil {
		return fmt.Errorf("error receiving message: %v", err)
	}

	return nil
}

func convertToTGXModel(msg *pubsub.Message) *models.MessageModel {
	s := models.PubSubMessageModel{}
	err := json.Unmarshal(msg.Data, &s)
	if err != nil {
		return &models.MessageModel{}
	}

	attributes := models.MessageAttributes{}
	err = mapstructure.Decode(msg.Attributes, &attributes)
	if err != nil {
		return &models.MessageModel{}
	}

	return &models.MessageModel{
		Payload:         s.Payload,
		SpecificPayload: s.SpecificPayload,
		Attributes:      attributes,
	}
}
