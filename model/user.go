package model

import (
	"time"
)

/******sql******
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(50) DEFAULT NULL COMMENT 'UUID',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除 0-未删除 1-已删除',
  `name` varchar(255) DEFAULT NULL COMMENT '名称',
  `age` int DEFAULT NULL COMMENT '年龄',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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