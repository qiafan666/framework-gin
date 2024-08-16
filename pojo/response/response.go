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
	UserList       []User `json:"user_list"`
}

// User User列表单项内容
type User struct {
	UUID        string    `json:"uuid"` // UUID
	CreatedTime time.Time `json:"created_time"`
	Name        string    `json:"name"` // 名称
	Age         int       `json:"age"`  // 年龄
}

// ================================================================================
// --------------------------------user_version表---------------------------------
// ================================================================================

// UserVersionCreate UserVersion表创建返回参数
type UserVersionCreate struct{}

// UserVersionDelete UserVersion表删除返回参数
type UserVersionDelete struct {
}

// UserVersionUpdate UserVersion表更新返回参数
type UserVersionUpdate struct {
}

// UserVersionList UserVersion表列表返回参数
type UserVersionList struct {
	BasePagination  `json:"-"`
	UserVersionList []UserVersion `json:"user_version_list"`
}

// UserVersion UserVersion列表单项内容
type UserVersion struct {
	UUID        string    `json:"uuid"` // UUID
	CreatedTime time.Time `json:"created_time"`
	UserID      int64     `json:"user_id"`    // user主键ID
	VersionID   int64     `json:"version_id"` // version主键ID
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
	VersionList    []Version `json:"version_list"`
}

// Version Version列表单项内容
type Version struct {
	UUID        string    `json:"uuid"` // UUID
	CreatedTime time.Time `json:"created_time"`
	Version     string    `json:"version"` // 版本号
}
