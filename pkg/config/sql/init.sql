CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id，全局唯一',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名，不可重复',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码，采用md5加密',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'gorm维护，创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'gorm维护，更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'gorm维护，删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user` (`username`) USING HASH COMMENT '用户名唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';