package mq

import "log"

var done chan bool

// StartConsume : Receive messages
func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	// Consume messages from the channel
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  // Auto acknowledge
		false, // Non-unique consumer
		false, // RabbitMQ can only be set to false
		false, // noWait, false means it will block until a message arrives
		nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	done = make(chan bool)

	go func() {
		// Loop to read data from the channel
		for d := range msgs {
			processErr := callback(d.Body)
			if processErr {
				// TODO: Write the task to the error queue for later processing
			}
		}
	}()

	// Receive the signal from done, will block until a message arrives to avoid exiting this function
	<-done

	// Close the channel
	channel.Close()
}

// StopConsume : Stop listening to the queue
func StopConsume() {
	done <- true
}
