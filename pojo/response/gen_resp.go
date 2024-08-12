package response

import (
	"time"
)

// ================================================================================
// ------------------------------------user表-------------------------------------
// ================================================================================

// UserCreate User表创建返回参数
type UserCreate struct{}

// UserDelete User表删除返回参数
type UserDelete struct {
}

// UserUpdate User表更新返回参数
type UserUpdate struct {
}

// UserList User表列表返回参数
type UserList struct {
	BasePagination `json:"-"`
	ID             int64     `json:"id"`   // 主键ID
	UUID           string    `json:"uuid"` // UUID
	CreatedTime    time.Time `json:"created_time"`
	Name           string    `json:"name"` // 名称
	Age            int       `json:"age"`  // 年龄
}

// ================================================================================
// -----------------------------------version表-----------------------------------
// ================================================================================

// VersionCreate Version表创建返回参数
type VersionCreate struct{}

// VersionDelete Version表删除返回参数
type VersionDelete struct {
}

// VersionUpdate Version表更新返回参数
type VersionUpdate struct {
}

// VersionList Version表列表返回参数
type VersionList struct {
	BasePagination `json:"-"`
	ID             int64     `json:"id"`   // 主键ID
	UUID           string    `json:"uuid"` // UUID
	CreatedTime    time.Time `json:"created_time"`
	Version        string    `json:"version"` // 版本号
}
