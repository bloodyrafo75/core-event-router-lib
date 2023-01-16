package pubsubService

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type PubSubConnection struct {
	ctx         context.Context
	project_id  string
	topic_name  string
	credentials string
}

func NewPubSubConnection(ctx context.Context, project_id string, topic_name string, credentials string) *PubSubConnection {
	return &PubSubConnection{
		ctx:         ctx,
		project_id:  project_id,
		topic_name:  topic_name,
		credentials: credentials,
	}
}

func (c *PubSubConnection) Connect() (*pubsub.Client, error) {
	// TODO set connection timeout
	creds, err := google.CredentialsFromJSON(c.ctx, []byte(c.credentials), pubsub.ScopePubSub)
	if err != nil {
		return nil, err
	}

	client, err := pubsub.NewClient(c.ctx, c.project_id, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *PubSubConnection) GetTopic(client pubsub.Client) (*pubsub.Topic, error) {
	topic := client.Topic(c.topic_name)
	if topic == nil {
		return nil, fmt.Errorf("topic %s doesn't exist", c.topic_name)
	}
	return topic, nil
}
