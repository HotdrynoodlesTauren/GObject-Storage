package mq

import (
	cmn "gobject-storage/common"
)

// TransferData : Struct for data to be sent to RabbitMQ
type TransferData struct {
	FileHash      string        // The hash of the file
	CurLocation   string        // Current location of the file
	DestLocation  string        // Destination location for the file
	DestStoreType cmn.StoreType // Destination storage type for the file
}
