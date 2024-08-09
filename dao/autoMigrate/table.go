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
	&Test{},
}

// BaseModel is base struct for entity, copy from gorm.Model
type BaseModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"` // 主键ID
	UUID        string    `gorm:"column:uuid;type:varchar(50);comment:UUID" json:"uuid"`     // UUID
	CreatedTime time.Time `gorm:"column:created_time;type:datetime(3);default:CURRENT_TIMESTAMP(3)" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime(3);default:CURRENT_TIMESTAMP(3) on update CURRENT_TIMESTAMP(3)" json:"updated_time"`
	IsDeleted   int8      `gorm:"column:is_deleted;default:0;comment:是否删除 0-未删除 1-已删除" json:"is_deleted"` // 是否删除 0-未删除 1-已删除
}

type Test struct {
	BaseModel
	Name string `gorm:"column:name;type:varchar(255);comment:名称" json:"name"`
	Age  int32  `gorm:"column:age;comment:年龄" json:"age"`
	Desc string `gorm:"column:desc;type:text;comment:描述" json:"desc"`
	Src  string `gorm:"column:src;type:varchar(255);comment:来源" json:"src"`
}

// TableName .
func (s *Test) TableName() string {
	return "test"
}
