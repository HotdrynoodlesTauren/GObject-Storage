package mq

import (
	"gobject-storage/config"

	"github.com/streadway/amqp"
)

// Publish : Publishes a message
func Publish(exchange, routingKey string, msg []byte) bool {
	// Initialize the channel
	if !initChannel(config.RabbitURL) {
		return false
	}

	// Publish the message to the channel
	if nil == channel.Publish(
		exchange,
		routingKey,
		false, // If there is no corresponding queue, the message will be discarded
		false, //
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg}) {
		return true
	}
	return false
}
