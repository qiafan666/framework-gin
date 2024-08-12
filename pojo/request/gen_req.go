package request

import ()

// ================================================================================
// ------------------------------------user表-------------------------------------
// ================================================================================

// UserCreate User表创建请求参数
type UserCreate struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	Name             string `json:"name"` // 名称
	Age              int    `json:"age"`  // 年龄
}

// UserDelete User表删除请求参数
type UserDelete struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	UUID             string `json:"uuid"` // UUID
}

// UserUpdate User表更新请求参数
type UserUpdate struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	UUID             *string `json:"uuid"` // UUID
	Name             *string `json:"name"` // 名称
	Age              *int    `json:"age"`  // 年龄
}

// UserList User表列表请求参数
type UserList struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	BasePagination
}

// ================================================================================
// --------------------------------user_version表---------------------------------
// ================================================================================

// UserVersionCreate UserVersion表创建请求参数
type UserVersionCreate struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	UserID           int64 `json:"user_id"`    // user主键ID
	VersionID        int64 `json:"version_id"` // version主键ID
}

// UserVersionDelete UserVersion表删除请求参数
type UserVersionDelete struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	UUID             string `json:"uuid"` // UUID
}

// UserVersionUpdate UserVersion表更新请求参数
type UserVersionUpdate struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	UUID             *string `json:"uuid"`       // UUID
	UserID           *int64  `json:"user_id"`    // user主键ID
	VersionID        *int64  `json:"version_id"` // version主键ID
}

// UserVersionList UserVersion表列表请求参数
type UserVersionList struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	BasePagination
}

// ================================================================================
// -----------------------------------version表-----------------------------------
// ================================================================================

// VersionCreate Version表创建请求参数
type VersionCreate struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	Version          string `json:"version"` // 版本号
}

// VersionDelete Version表删除请求参数
type VersionDelete struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	UUID             string `json:"uuid"` // UUID
}

// VersionUpdate Version表更新请求参数
type VersionUpdate struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	UUID             *string `json:"uuid"`    // UUID
	Version          *string `json:"version"` // 版本号
}

// VersionList Version表列表请求参数
type VersionList struct {
	BaseRequest      `json:"-"`
	BaseTokenRequest `json:"-"`
	BasePagination
}
