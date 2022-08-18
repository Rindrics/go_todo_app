CREATE TABLE `user`
(
	`id`		BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'User ID',
	`name`		VARCHAR(20) NOT NULL COMMENT 'User name',
	`password` 	VARCHAR(80) NOT NULL COMMENT 'Password hash',
	`role` 		VARCHAR(80) NOT NULL COMMENT 'Role',
	`created` 	DATETIME(6) NOT NULL COMMENT 'Date record created',
	`modified` 	DATETIME(6) NOT NULL COMMENT 'Date record updated',
        PRIMARY KEY (`id`),
        UNIQUE KEY `uix_name` (`name`) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User';

CREATE TABLE `task`
(
	`id`		BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Task ID',
	`title`		VARCHAR(128) NOT NULL COMMENT 'Task title',
	`status` 	VARCHAR(20) NOT NULL COMMENT 'Task status',
	`created` 	DATETIME(6) NOT NULL COMMENT 'Date record created',
	`modified` 	DATETIME(6) NOT NULL COMMENT 'Date record updated',
        PRIMARY KEY (`id`)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Task';
