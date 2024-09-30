package model

import (
	"time"
)

/******sql******
CREATE TABLE `user_version` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(50) DEFAULT NULL COMMENT 'UUID',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除 0-未删除 1-已删除',
  `user_id` bigint DEFAULT NULL COMMENT 'user主键ID',
  `version_id` bigint DEFAULT NULL COMMENT 'version主键ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_uuid` (`uuid`),
  KEY `idx_user_version_id` (`user_id`,`version_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
******sql******/
// UserVersion [...]
type UserVersion struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"-"` // 主键ID
	UUID        string    `gorm:"column:uuid" json:"uuid"`       // UUID
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
	IsDeleted   int8      `gorm:"column:is_deleted" json:"is_deleted"` // 是否删除 0-未删除 1-已删除
	UserID      int64     `gorm:"column:user_id" json:"user_id"`       // user主键ID
	VersionID   int64     `gorm:"column:version_id" json:"version_id"` // version主键ID
}

// TableName get sql table name.获取数据库表名
func (m *UserVersion) TableName() string {
	return "user_version"
}

// UserVersionColumns get sql column name.获取数据库列名
var UserVersionColumns = struct {
	ID          string
	UUID        string
	CreatedTime string
	UpdatedTime string
	IsDeleted   string
	UserID      string
	VersionID   string
}{
	ID:          "id",
	UUID:        "uuid",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
	IsDeleted:   "is_deleted",
	UserID:      "user_id",
	VersionID:   "version_id",
}
