# Core-event-router
Go library to publish/consume to/from Pubsub


For local testing: 

- Local docker to get a pubSub service.

```python
cd ./third_parties
docker-compose up -d
cd ..
```

- Set values for `configs/.env`

```python
EVENT_ROUTER_CLIENTID='IAM'
PUBSUB_HOST = "127.0.0.1:8262" # Only for working with local docker
PROJECT_ID = "my-project-id"
TOPIC_NAME = "topic-name"
PORT = "3000" # Only for working with local docker
SUBSCRIPTION_ID="example-subscription"
LOCAL_MODE="1"
CONSUMER_CREDENTIALS=""
PRODUCER_CREDENTIALS=""
```

- Maybe you will need to execute

```python
export GO111MODULE=on
```


- Start the receiver example (listener/worker)

```python
go run cmd/receiver/receiver.go &
```


- Start the sender example (listener/worker)

```python
go run cmd/sender/sender.go 
```