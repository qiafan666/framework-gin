package model

import (
	"time"
)

/******sql******
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(50) DEFAULT NULL COMMENT 'UUID',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `is_deleted` tinyint(4) DEFAULT '0' COMMENT '是否删除 0-未删除 1-已删除',
  `name` varchar(255) DEFAULT NULL COMMENT '名称',
  `age` int(11) DEFAULT NULL COMMENT '年龄',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
******sql******/
// User [...]
type User struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"-"` // 主键ID
	UUID        string    `gorm:"column:uuid" json:"uuid"`       // UUID
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
	IsDeleted   int8      `gorm:"column:is_deleted" json:"is_deleted"` // 是否删除 0-未删除 1-已删除
	Name        string    `gorm:"column:name" json:"name"`             // 名称
	Age         int       `gorm:"column:age" json:"age"`               // 年龄
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}

// UserColumns get sql column name.获取数据库列名
var UserColumns = struct {
	ID          string
	UUID        string
	CreatedTime string
	UpdatedTime string
	IsDeleted   string
	Name        string
	Age         string
}{
	ID:          "id",
	UUID:        "uuid",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
	IsDeleted:   "is_deleted",
	Name:        "name",
	Age:         "age",
}

/******sql******
CREATE TABLE `user_version` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(50) DEFAULT NULL COMMENT 'UUID',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `is_deleted` tinyint(4) DEFAULT '0' COMMENT '是否删除 0-未删除 1-已删除',
  `user_id` bigint(20) DEFAULT NULL COMMENT 'user主键ID',
  `version_id` bigint(20) DEFAULT NULL COMMENT 'version主键ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
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

/******sql******
CREATE TABLE `version` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(50) DEFAULT NULL COMMENT 'UUID',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `is_deleted` tinyint(4) DEFAULT '0' COMMENT '是否删除 0-未删除 1-已删除',
  `version` varchar(50) DEFAULT NULL COMMENT '版本号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_uuid` (`uuid`),
  UNIQUE KEY `uix_version` (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
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
