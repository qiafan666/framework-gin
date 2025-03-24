package autoMigrate

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func CreateTables(db *gorm.DB) {
	for _, table := range Tables {
		// 检查表是否存在
		if !db.Migrator().HasTable(table) {
			// 表不存在时创建
			err := db.Migrator().CreateTable(table)
			if err != nil {
				panic(err)
			}

			tableName := table.(TableInterface).TableName()
			comment := table.(TableInterface).Component()
			query := fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", tableName, comment)
			_ = db.Exec(query).Error
		}
	}
}

var Tables = []interface{}{
	&User{},
	&Version{},
	&UserVersion{},
}

type TableInterface interface {
	Component() string
	TableName() string
}

// BaseModel is base struct for entity, copy from gorm.Model
type BaseModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement;not null;comment:主键ID"`                // 主键ID
	UUID        string    `gorm:"column:uuid;type:varchar(50);not null;comment:UUID;uniqueIndex:uix_uuid"` // UUID
	CreatedTime time.Time `gorm:"column:created_time;type:timestamp(3);not null;default:CURRENT_TIMESTAMP(3);comment:创建时间"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:timestamp(3);not null;default:CURRENT_TIMESTAMP(3) on update CURRENT_TIMESTAMP(3);comment:更新时间"`
	IsDeleted   int8      `gorm:"column:is_deleted;type:tinyint;default:0;not null;comment:是否删除 0-未删除 1-已删除;index:idx_is_deleted"` // 是否删除 0-未删除 1-已删除
}

// -------------------- 用户表 ----------------------

type User struct {
	BaseModel
	Name string `gorm:"column:name;type:varchar(255);not null;comment:名称"`
	Age  int32  `gorm:"column:age;not null;comment:年龄"`
}

func (s *User) TableName() string {
	return "user"
}
func (s *User) Component() string {
	return "用户表"
}

// -------------------- 版本表 ----------------------

type Version struct {
	BaseModel
	Version string `gorm:"column:version;type:varchar(50);not null;comment:版本号;uniqueIndex:uix_version"`
}

func (s *Version) TableName() string {
	return "version"
}
func (s *Version) Component() string {
	return "版本表"
}

// -------------------- 用户版本关系表 ----------------------

type UserVersion struct {
	BaseModel
	UserID    int64 `gorm:"column:user_id;comment:user主键ID;not null;index:idx_user_version_id"`
	VersionID int64 `gorm:"column:version_id;comment:version主键ID;not null;index:idx_user_version_id"`
}

func (s *UserVersion) TableName() string {
	return "user_version"
}
func (s *UserVersion) Component() string {
	return "用户版本关系表"
}
