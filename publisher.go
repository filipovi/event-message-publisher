package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/filipovi/event-message-publisher/config"
	"github.com/filipovi/event-message-publisher/rabbitmq"
)

/**
 * @TODO
 *  1. add a config file
 *  2. cli parameters
 */

// EventMessage represents a Message sent in RabbitMQ
type EventMessage struct {
	ID        string              `json:"id"`
	ModelType string              `json:"model_type"`
	Event     string              `json:"event"`
	Metadata  map[string][]string `json:"metadata"`
}

// Env is a structure who contains the Rabbitmq channel
type Env struct {
	channel rabbitmq.Channel
}

func failOnError(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", msg, err)
	panic(fmt.Sprintf("%s: %s", msg, err))
}

func createMessage() EventMessage {
	message := EventMessage{
		ID:        "cfb83af8-99d1-11e6-9f33-a24fc0d9649c",
		Event:     "updated",
		ModelType: "user",
		Metadata: map[string][]string{
			"updated_fields": []string{
				"lastLogin",
				"updatedAt",
			},
		},
	}
	return message
}

func connect(cfg config.Config) (*Env, error) {
	rabbitmq, err := rabbitmq.New(cfg.Rabbitmq.URL)
	if nil != err {
		return nil, err
	}
	log.Println("Rabbitmq connected!")

	env := &Env{
		channel: *rabbitmq,
	}

	return env, nil
}

func main() {
	cfg, err := config.New("config.json")
	failOnError(err, "Failed to read config.json")

	env, err := connect(cfg)
	failOnError(err, "Failed to get a apmq channel")

	err = env.channel.NewExchange(rabbitmq.Exchange{
		Name:         "event_message.mailchimp",
		ExchangeType: "fanout",
		Durable:      true,
	})
	failOnError(err, "Failed to create an apmq exchange")

	message := createMessage()
	body, _ := json.Marshal(&message)

	err = env.channel.Send(body, "event_message")
	failOnError(err, "Failed to send a message")

	log.Printf(" [x] Sent %s", body)
}
