package config

const (
	AsyncTransferEnable  = true
	TransExchangeName    = "uploadserver.trans"
	TransOSSQueueName    = "uploadserver.trans.oss"
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	TransOSSRoutingKey   = "oss"
)

var (
	RabbitURL = "amqp://guest@127.0.0.1:5672/"
)
