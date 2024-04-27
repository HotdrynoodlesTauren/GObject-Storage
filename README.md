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

## Preferred Environment Settings
Mac OS Version 13+

Docker version 25.0.3

Go Version 1.21.1

## Prerequisites
### 1. Go Language
1. Make sure that you have Go Language installed.
2. Install all the packages needed for the project by running the following command in the project folder:
```
go mod tidy
```

### 2. Setup MySQL Server
1. Prepare a MySQL 5.7 Instance within a Docker container. If desired, set up a lead-follower based cluster.
2. Creat the target database and tables in the Mysql Instance.
```
create database <database name> default character set utf8;
use <database name>;

CREATE TABLE `tbl_file` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `file_sha1` CHAR(40) NOT NULL DEFAULT '' COMMENT 'File hash',
    `file_name` VARCHAR(256) NOT NULL DEFAULT '' COMMENT 'File name',
    `file_size` BIGINT(20) DEFAULT 0 COMMENT 'File size',
    `file_addr` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT 'File storage location',
    `created_at` DATETIME DEFAULT NOW() COMMENT 'Creation date',
    `updated_at` DATETIME DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP() COMMENT 'Update date',
    `status` INT(11) NOT NULL DEFAULT 0 COMMENT 'Status (available/disabled/deleted)',
    `ext1` INT(11) DEFAULT 0 COMMENT 'Reserved field 1',
    `ext2` TEXT COMMENT 'Reserved field 2',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_hash` (`file_sha1`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tbl_user` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `user_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Username',
    `email` VARCHAR(64) DEFAULT '' COMMENT 'Email address',
    `phone` VARCHAR(128) DEFAULT '' COMMENT 'Phone number',
    `email_validated` TINYINT(1) DEFAULT 0 COMMENT 'Whether the email is validated',
    `phone_validated` TINYINT(1) DEFAULT 0 COMMENT 'Whether the phone is validated',
    `signup_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Registration date',
    `last_active` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last active time',
    `profile` TEXT COMMENT 'User profile',
    `status` INT(11) NOT NULL DEFAULT 0 COMMENT 'Account status (enabled/disabled/locked/deleted flag, etc.)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_phone` (`phone`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `tbl_user_token` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `user_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Username',
    `user_token` CHAR(40) NOT NULL DEFAULT '' COMMENT 'User login token',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `tbl_user_file` (
    `id` INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_name` VARCHAR(64) NOT NULL,
    `file_shal` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'File hash',
    `file_size` BIGINT(20) DEFAULT '0' COMMENT 'File size in bytes',
    `file_name` VARCHAR(256) NOT NULL DEFAULT '' COMMENT 'File name',
    `upload_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Upload time',
    `last_update` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last modification time',
    `status` INT(11) NOT NULL DEFAULT '0' COMMENT 'File status (0: normal, 1: deleted, 2: disabled)',
    UNIQUE KEY `idx_user_file` (`user_name`, `file_shal`),
    KEY `idx_status` (`status`),
    KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
alter table tbl_user_file drop index idx_user_file;
```

3. Find `db/mysql/conn.go` and modify the following lines according to your MySQL setup:
```
db, _ = sql.Open("mysql", "<root username>:<root pwd>@tcp(127.0.0.1:<port used for Mysql>)/<database name>?charset=utf8")
```

### 3. Setup Redis Server
1. Prepare a Redis 5.0.3 server within a Docker container.
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
1. Prepare for a RabbitMQ server within a Docker container.
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


