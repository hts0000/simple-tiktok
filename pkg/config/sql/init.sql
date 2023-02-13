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

CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评论id，全局唯一',
  `user_id` bigint unsigned NOT NULL COMMENT '评论者id',
  `video_id` bigint unsigned NOT NULL COMMENT '评论视频id',
  `content` varchar(500) NOT NULL COMMENT '评论内容',
  `create_date` varchar(15) NOT NULL COMMENT '评论创建日期',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'gorm维护，创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'gorm维护，更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'gorm维护，删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_video` (`video_id`) USING BTREE COMMENT '视频id索引',
  FOREIGN KEY (`user_id`) references user(`id`),
  FOREIGN KEY (`video_id`) references videos(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='评论表';

CREATE TABLE `likes` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id` bigint unsigned NOT NULL COMMENT '点赞用户id',
    `video_id` bigint unsigned NOT NULL COMMENT '被点赞的视频id',
    `cancel` tinyint(4) NOT NULL DEFAULT '0' COMMENT '默认点赞为0，取消赞为1',
    PRIMARY KEY (`id`),
    UNIQUE KEY `userIdtoVideoIdIdx` (`user_id`,`video_id`) USING BTREE,
    KEY `userIdIdx` (`user_id`) USING BTREE,
    KEY `videoIdx` (`video_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='点赞表';

