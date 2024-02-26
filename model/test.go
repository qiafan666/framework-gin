package model

import (
	"time"
)

/******sql******
CREATE TABLE `test` (
  `id` bigint(255) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `test` varchar(191) NOT NULL,
  `is_delete` tinyint(2) NOT NULL DEFAULT '2' COMMENT '1:yes 2:no',
  `created_time` timestamp NULL DEFAULT NULL,
  `updated_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4
******sql******/
// Version [...]
type Version struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"-"` // 主键id
	Test        string    `gorm:"column:test" json:"test"`
	IsDelete    int8      `gorm:"column:is_delete" json:"is_delete"` // 1:yes 2:no
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
}

// TableName get sql table name.获取数据库表名
func (m *Version) TableName() string {
	return "test"
}

// VersionColumns get sql column name.获取数据库列名
var VersionColumns = struct {
	ID          string
	Test        string
	IsDelete    string
	CreatedTime string
	UpdatedTime string
}{
	ID:          "id",
	Test:        "test",
	IsDelete:    "is_delete",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
}
