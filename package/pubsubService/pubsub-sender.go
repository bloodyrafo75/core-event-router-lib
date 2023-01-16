package pubsubService

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/fatih/structs"
	"gitlab.xmltravelgate.com/core/core-event-router/models"
)

type PubSubSender struct {
	ctx        context.Context
	connection PubSubConnection
}

func NewPubSubSender(ctx context.Context, connection PubSubConnection) *PubSubSender {
	return &PubSubSender{
		ctx:        ctx,
		connection: connection,
	}
}

func (s *PubSubSender) send(payload string, attributes map[string]string) (serverID *string, err error) {
	client, err := s.connection.Connect()
	if err != nil {
		return nil, err
	}
	topic, err := s.connection.GetTopic(*client)
	if err != nil {
		return nil, err
	}
	res := topic.Publish(s.ctx, &pubsub.Message{Attributes: attributes, Data: []byte(payload)})
	messageId, err := res.Get(s.ctx)
	if err != nil {
		_ = close(*client)
		return nil, err
	}
	_ = close(*client)

	return &messageId, nil
}

func (s *PubSubSender) Publish(msg models.MessageModel, clientID string) (*models.Response, error) {
	msg.Attributes.Src = clientID
	attributes := getAttributes(msg.Attributes)

	pubsubMsg := models.PubSubMessageModel{
		Payload:         msg.Payload,
		SpecificPayload: msg.SpecificPayload,
	}

	eventStr, err := json.Marshal(pubsubMsg)
	if err != nil {
		return nil, err
	}

	var rsp models.Response
	rsp.MessageID, err = s.send(string(eventStr), attributes)
	if err != nil {
		return nil, err
	}

	return &rsp, nil
}

func getAttributes(eventAttributes models.MessageAttributes) map[string]string {
	attributes := make(map[string]string)

	keys := structs.Names(eventAttributes)
	values := structs.Values(eventAttributes)

	for i, key := range keys {
		attributes[strings.ToLower(key)] = fmt.Sprintf("%v", values[i])
	}

	return attributes
}

func close(client pubsub.Client) error {
	err := client.Close()
	if err != nil {
		return err
	}
	return nil

}
