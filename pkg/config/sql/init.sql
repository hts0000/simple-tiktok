CREATE DATABASE `tiktok` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_general_ci';
use tiktok;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id，全局唯一',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名，不可重复',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码，采用md5加密',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'gorm维护，创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'gorm维护，更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'gorm维护，删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`) COMMENT '用户名唯一索引',
  KEY `idx_delete_at` (`deleted_at`) USING BTREE COMMENT 'gorm删除时间索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

CREATE TABLE `follow` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint unsigned NOT NULL COMMENT '用户id，全局唯一',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名，不可重复',
  `follower_id` bigint unsigned NOT NULL COMMENT '粉丝id',
  `follower_name` varchar(255) NOT NULL DEFAULT '' COMMENT '粉丝名，不可重复',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'gorm维护，创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'gorm维护，更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'gorm维护，删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`) USING BTREE COMMENT '用户id索引',
  KEY `idx_follower_id` (`follower_id`) USING BTREE COMMENT '粉丝id索引',
  KEY `idx_delete_at` (`deleted_at`) USING BTREE COMMENT 'gorm删除时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户关系表';

CREATE TABLE `videos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '视频id，全局唯一',
  `user_id` bigint unsigned NOT NULL COMMENT '上传者id',
  `title` varchar(50) NOT NULL COMMENT '视频标题',
  `type` varchar(10) NOT NULL COMMENT '视频类型',
  `favorite_count` int unsigned DEFAULT 0 NOT NULL COMMENT '点赞数量',
  `comment_count` int unsigned DEFAULT 0 NOT NULL COMMENT '评论数量',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'gorm维护，创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'gorm维护，更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'gorm维护，删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`) COMMENT '用户id索引',
  KEY `idx_created_at` (`created_at`) USING BTREE COMMENT 'gorm创建时间索引',
  FOREIGN KEY (`user_id`) references user(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='视频表';