CREATE TABLE `tbl_file` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT 'File hash',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'File name',
    `file_size` bigint(20) DEFAULT 0 COMMENT 'File size',
    `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT 'File storage location',
    `created_at` datetime DEFAULT NOW() COMMENT 'Creation date',
    `updated_at` datetime DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP() COMMENT 'Update date',
    `status` int(11) NOT NULL DEFAULT 0 COMMENT 'Status (available/disabled/deleted)',
    `ext1` int(11) DEFAULT 0 COMMENT 'Reserved field 1',
    `ext2` text COMMENT 'Reserved field 2',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_hash` (`file_sha1`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tbl_user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'Username',
    `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT 'Encoded password for the user',
    `email` varchar(64) DEFAULT '' COMMENT 'Email address',
    `phone` varchar(128) DEFAULT '' COMMENT 'Phone number',
    `email_validated` tinyint(1) DEFAULT 0 COMMENT 'Whether the email is validated',
    `phone_validated` tinyint(1) DEFAULT 0 COMMENT 'Whether the phone is validated',
    `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Registration date',
    `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last active time',
    `profile` text COMMENT 'User profile',
    `status` int(11) NOT NULL DEFAULT 0 COMMENT 'Account status (enabled/disabled/locked/deleted flag, etc.)',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

