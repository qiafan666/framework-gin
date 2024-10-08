package model

import (
	"time"
)

/******sql******
CREATE TABLE `version` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(50) NOT NULL COMMENT 'UUID',
  `created_time` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `is_deleted` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 0-未删除 1-已删除',
  `version` varchar(50) NOT NULL COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_uuid` (`uuid`),
  UNIQUE KEY `uix_version` (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
******sql******/
// Version [...]
type Version struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"-"` // 主键ID
	UUID        string    `gorm:"column:uuid" json:"uuid"`       // UUID
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
	IsDeleted   int8      `gorm:"column:is_deleted" json:"is_deleted"` // 是否删除 0-未删除 1-已删除
	Version     string    `gorm:"column:version" json:"version"`       // 版本号
}

// TableName get sql table name.获取数据库表名
func (m *Version) TableName() string {
	return "version"
}

// VersionColumns get sql column name.获取数据库列名
var VersionColumns = struct {
	ID          string
	UUID        string
	CreatedTime string
	UpdatedTime string
	IsDeleted   string
	Version     string
}{
	ID:          "id",
	UUID:        "uuid",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
	IsDeleted:   "is_deleted",
	Version:     "version",
}
