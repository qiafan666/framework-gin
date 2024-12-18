package autoMigrate

import (
	"gorm.io/gorm"
	"time"
)

func CreateTables(db *gorm.DB) {
	for _, table := range Tables {
		err := db.Migrator().DropTable(table)
		if err != nil {
			panic(err)
		}
		err = db.Migrator().CreateTable(table)
		if err != nil {
			panic(err)
		}
	}
}

var Tables = []interface{}{
	&User{},
	&Version{},
	&UserVersion{},
}

// BaseModel is base struct for entity, copy from gorm.Model
type BaseModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID"`                         // 主键ID
	UUID        string    `gorm:"column:uuid;type:varchar(50);not null;comment:UUID;uniqueIndex:uix_uuid"` // UUID
	CreatedTime time.Time `gorm:"column:created_time;type:timestamp(3);not null;default:CURRENT_TIMESTAMP(3)"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:timestamp(3);not null;default:CURRENT_TIMESTAMP(3) on update CURRENT_TIMESTAMP(3)"`
	IsDeleted   int8      `gorm:"column:is_deleted;default:0;not null;comment:是否删除 0-未删除 1-已删除"` // 是否删除 0-未删除 1-已删除
}

// User 用户表
type User struct {
	BaseModel
	Name string `gorm:"column:name;type:varchar(255);not null;comment:名称"`
	Age  int32  `gorm:"column:age;not null;comment:年龄"`
}

func (s *User) TableName() string {
	return "user"
}

// Version 版本表
type Version struct {
	BaseModel
	Version string `gorm:"column:version;type:varchar(50);not null;comment:版本号;uniqueIndex:uix_version"`
}

func (s *Version) TableName() string {
	return "version"
}

// UserVersion 用户版本关系表
type UserVersion struct {
	BaseModel
	UserID    int64 `gorm:"column:user_id;comment:user主键ID;not null;index:idx_user_version_id"`
	VersionID int64 `gorm:"column:version_id;comment:version主键ID;not null;index:idx_user_version_id"`
}

func (s *UserVersion) TableName() string {
	return "user_version"
}
