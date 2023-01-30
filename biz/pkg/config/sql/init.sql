CREATE TABLE tiktok.`user` (
	id BIGINT(20) unsigned auto_increment NOT NULL COMMENT '用户id，全局唯一',
	name varchar(255) NOT NULL COMMENT '用户名，不可重复',
	password varchar(255) NOT NULL COMMENT '用户密码，采用md5加密',
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'gorm维护，创建时间',
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'gorm维护，更新时间',
    deleted_at timestamp NULL DEFAULT NULL COMMENT 'gorm维护，删除时间',
	CONSTRAINT id PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_general_ci
COMMENT='用户表'
AUTO_INCREMENT=1;

CREATE INDEX user_name_IDX USING BTREE ON tiktok.`user` (name,password);
