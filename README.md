# GObject-Storage: End-to-end Scaclable Distributed Object Storage

## Contributors
Yi Zhou, Yachen Yin

## Introduction
Leveraging the power of Golang and industry-leading technologies such as `MySQL`, `Redis`, `Aliyun` and `RabbitMQ`, this project aims to deliver a robust solution for storing and managing user files while ensuring high performance, fault tolerance, and data integrity.

## Key Features:
#### Distributed Architecture: 
GObjectStorage utilizes a distributed architecture to ensure scalability and fault tolerance.

#### Object Storage:
Users can store and retrieve objects of any size, including files, documents, images, and multimedia content.

#### Access Control:
User metadata stored in MySQL enables secure and personalized storage by linking each user's storage space with their profile。

#### Duplicate Item Detection:
Utilizing file metadata stored in `MySQL`, the system optimizes storage efficiency through hash algorithms designed to identify duplicate items.

#### Block Metadata Cache:
`Redis` stores block-level metadata, enabling chunked and resumable uploads with efficient management and seamless file upload resumption.

#### Message Queue:
`RabbitMQ` facilitates communication and coordination between system components, ensuring reliable message delivery and enabling asynchronous processing of storage operations.

#### Scalability and Reliability:
Leveraging `Aliyun OSS`'s distributed storage capabilities, our system achieves automatic load balancing, data redundancy, fault tolerance, and high availability through efficient replication and data striping.


## Prerequisites
### 1. Go Language
1. Make sure that you have Go Language installed.
2. Install all the packages needed for the project by running the following command in the project folder:
```
go mod tidy
```

### 2. Setup MySQL Server
1. Prepare a MySQL Instance. If desired, set up a lead-follower based cluster.
2. Find `db/mysql/conn.go` and modify the following line according to your MySQL setup:
```
db, _ = sql.Open("mysql", "<root username>:<root pwd>@tcp(127.0.0.1:<port used for Mysql>)/<database name>?charset=utf8")
```

### 3. Setup Redis Server
1. Prepare a Redis server.
2. Find `cache/redis/conn.go` and modify the following line according to your Redis setup:
```
redisHost = "127.0.0.1:<port used for Redis>"
```

### 3. Setup Aliyun OSS
1. Sign up or Log in to your Aliyun account.
2. Set up an OSS instance and create a bucket that will be used for the storage.
3. Create 'config/oss.go' and add the following code to your file based on your Aliyun setup
```
package config

const (
	OSSBucket          = "<Your bucket name>"
	OSSEndpoint        = "<Your OSS end point>"
	OSSAccesskeyID     = "<Your access key>"
	OSSAccessKeySecret = "<Your access key secret>"
)
```
**Do not share your __OSSAccesskeyID__ and __OSSAccessKeySecret__ with others!**

### 4. Setup RabbitMQ
1. Prepare for a RabbitMQ server.
2. Generate a `Exchange` service and a message `Queue`.
3. Find `config/rabbit.go` and modify the following line according to your RabbitMQ setup:
```
TransExchangeName    = "<Your RabbitMQ Exhange Service name>"
TransOSSQueueName    = "<Your RabbitMQ Queue name>"
TransOSSErrQueueName = "<Your RabbitMQ Queue name>.err"
TransOSSRoutingKey   = "<Your RabbitMQ RoutingKey name>"

RabbitURL = "amqp://guest@127.0.0.1:<port used for RabbitMQ>/"
```

## Run the Program
1. Run the following command in your project folder:
```
go run ./service/upload/main.go
go run ./service/transfer/main.go
```
2. Navigate to `http://localhost:8080/user/signup` in a web browser.


