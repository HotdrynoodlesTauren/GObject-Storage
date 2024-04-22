package mq

import (
	"gobject-storage/config"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

// Receives notifications if the connection is unexpectedly closed
var notifyClose chan *amqp.Error

// UpdateRabbitHost : Update MQ host
func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}

// Init : Initialize MQ connection information
func Init() {
	// Initialize RabbitMQ connection only if asynchronous transfer is enabled
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}
	// Automatic reconnection on disconnection
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}

func initChannel(rabbitHost string) bool {
	if channel != nil {
		return true
	}

	// Get a rabbitmq connection
	var err error
	conn, err = amqp.Dial(rabbitHost)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Open a channel for publishing and receiving messages
	channel, err = conn.Channel()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
