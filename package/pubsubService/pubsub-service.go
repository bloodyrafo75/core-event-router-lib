package pubsubService

import (
	"context"

	"github.com/bloodyrafo75/core-event-router-lib/models"
)

var PubsubClient PubSubService

type Callback func(*models.MessageModel) error

type PubSubService struct {
	ctx        context.Context
	connection PubSubConnection
	receiver   PubSubReceiver
	sender     PubSubSender

	projectID   string
	topicID     string
	credentials string
	clientID    string
}

func NewPubSubService(ctx context.Context, projectID string, topicID string, credentials string, clientID string) *PubSubService {
	PubsubClient = PubSubService{
		ctx:         ctx,
		projectID:   projectID,
		topicID:     topicID,
		credentials: credentials,
		clientID:    clientID,
	}

	return &PubsubClient
}

func (s *PubSubService) StartConsumer(callback *Callback, subscriptionID string) error {

	s.connection = *NewPubSubConnection(s.ctx, s.projectID, s.topicID, s.credentials)
	s.receiver = *NewPubSubReceiver(s.ctx, s.connection, subscriptionID, *callback)

	err := s.receiver.Start()
	if err != nil {
		return err
	}
	return nil
}

func (s *PubSubService) NotifyEvent(msg *models.MessageModel) (*models.Response, error) {
	s.connection = *NewPubSubConnection(s.ctx, s.projectID, s.topicID, s.credentials)
	if s.sender == (PubSubSender{}) {
		s.sender = *NewPubSubSender(s.ctx, s.connection)
	}
	return s.sender.Publish(*msg, s.clientID)

}

func (s *PubSubService) CreatePubsubMsg(src string, prod string, _type string, stype string, operation string, payload string, specific_payload string) models.MessageModel {
	attr := models.MessageAttributes{
		Src:   src,
		Prod:  prod,
		Type:  _type,
		Stype: stype,
		Op:    operation,
	}
	return models.MessageModel{
		Payload:         payload,
		SpecificPayload: specific_payload,
		Attributes:      attr,
	}
}
