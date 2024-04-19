package main

import (
	"encoding/json"
	"log"
	"os"

	"gobject-storage/config"
	dblayer "gobject-storage/db"
	"gobject-storage/mq"
	"gobject-storage/store/oss"
)

// Transfer : Handles file transfer
func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	// Parse the message into TransferData struct
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Open the file for reading
	_, err = os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Upload the file to the destination location in OSS
	err = oss.Bucket().PutObjectFromFile(
		pubData.DestLocation,
		pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Update the file location in the database
	suc := dblayer.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	if !suc {
		return false
	}
	return true
}

func main() {
	log.Println("Start listening to the transfer message queue")
	// Start consuming messages from the transfer OSS queue
	mq.Init()
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}
