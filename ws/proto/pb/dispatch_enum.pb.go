// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v4.25.3
// source: dispatch_enum.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Grp int32

const (
	Grp_Def   Grp = 0
	Grp_Sys   Grp = 1
	Grp_Logic Grp = 2
)

// Enum value maps for Grp.
var (
	Grp_name = map[int32]string{
		0: "Def",
		1: "Sys",
		2: "Logic",
	}
	Grp_value = map[string]int32{
		"Def":   0,
		"Sys":   1,
		"Logic": 2,
	}
)

func (x Grp) Enum() *Grp {
	p := new(Grp)
	*p = x
	return p
}

func (x Grp) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Grp) Descriptor() protoreflect.EnumDescriptor {
	return file_dispatch_enum_proto_enumTypes[0].Descriptor()
}

func (Grp) Type() protoreflect.EnumType {
	return &file_dispatch_enum_proto_enumTypes[0]
}

func (x Grp) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Grp.Descriptor instead.
func (Grp) EnumDescriptor() ([]byte, []int) {
	return file_dispatch_enum_proto_rawDescGZIP(), []int{0}
}

type CmdSys int32

const (
	CmdSys_DefCmdSys           CmdSys = 0
	CmdSys_SubscribeOnlineUser CmdSys = 1
	CmdSys_KickOnlineUser      CmdSys = 2
	CmdSys_PushMessage         CmdSys = 3
)

// Enum value maps for CmdSys.
var (
	CmdSys_name = map[int32]string{
		0: "DefCmdSys",
		1: "SubscribeOnlineUser",
		2: "KickOnlineUser",
		3: "PushMessage",
	}
	CmdSys_value = map[string]int32{
		"DefCmdSys":           0,
		"SubscribeOnlineUser": 1,
		"KickOnlineUser":      2,
		"PushMessage":         3,
	}
)

func (x CmdSys) Enum() *CmdSys {
	p := new(CmdSys)
	*p = x
	return p
}

func (x CmdSys) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CmdSys) Descriptor() protoreflect.EnumDescriptor {
	return file_dispatch_enum_proto_enumTypes[1].Descriptor()
}

func (CmdSys) Type() protoreflect.EnumType {
	return &file_dispatch_enum_proto_enumTypes[1]
}

func (x CmdSys) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CmdSys.Descriptor instead.
func (CmdSys) EnumDescriptor() ([]byte, []int) {
	return file_dispatch_enum_proto_rawDescGZIP(), []int{1}
}

type CmdLogic int32

const (
	CmdLogic_DefCmdLogic CmdLogic = 0
	CmdLogic_Health      CmdLogic = 1
)

// Enum value maps for CmdLogic.
var (
	CmdLogic_name = map[int32]string{
		0: "DefCmdLogic",
		1: "Health",
	}
	CmdLogic_value = map[string]int32{
		"DefCmdLogic": 0,
		"Health":      1,
	}
)

func (x CmdLogic) Enum() *CmdLogic {
	p := new(CmdLogic)
	*p = x
	return p
}

func (x CmdLogic) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CmdLogic) Descriptor() protoreflect.EnumDescriptor {
	return file_dispatch_enum_proto_enumTypes[2].Descriptor()
}

func (CmdLogic) Type() protoreflect.EnumType {
	return &file_dispatch_enum_proto_enumTypes[2]
}

func (x CmdLogic) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CmdLogic.Descriptor instead.
func (CmdLogic) EnumDescriptor() ([]byte, []int) {
	return file_dispatch_enum_proto_rawDescGZIP(), []int{2}
}

type KickReason int32

const (
	KickReason_DefKickReason KickReason = 0
	// 只能一个端登录
	KickReason_OnlyOneClient KickReason = 1
)

// Enum value maps for KickReason.
var (
	KickReason_name = map[int32]string{
		0: "DefKickReason",
		1: "OnlyOneClient",
	}
	KickReason_value = map[string]int32{
		"DefKickReason": 0,
		"OnlyOneClient": 1,
	}
)

func (x KickReason) Enum() *KickReason {
	p := new(KickReason)
	*p = x
	return p
}

func (x KickReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (KickReason) Descriptor() protoreflect.EnumDescriptor {
	return file_dispatch_enum_proto_enumTypes[3].Descriptor()
}

func (KickReason) Type() protoreflect.EnumType {
	return &file_dispatch_enum_proto_enumTypes[3]
}

func (x KickReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use KickReason.Descriptor instead.
func (KickReason) EnumDescriptor() ([]byte, []int) {
	return file_dispatch_enum_proto_rawDescGZIP(), []int{3}
}

var File_dispatch_enum_proto protoreflect.FileDescriptor

var file_dispatch_enum_proto_rawDesc = []byte{
	0x0a, 0x13, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x2a, 0x22, 0x0a, 0x03, 0x47, 0x72, 0x70,
	0x12, 0x07, 0x0a, 0x03, 0x44, 0x65, 0x66, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x79, 0x73,
	0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x10, 0x02, 0x2a, 0x55, 0x0a,
	0x06, 0x43, 0x6d, 0x64, 0x53, 0x79, 0x73, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x65, 0x66, 0x43, 0x6d,
	0x64, 0x53, 0x79, 0x73, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65, 0x72, 0x10, 0x01, 0x12,
	0x12, 0x0a, 0x0e, 0x4b, 0x69, 0x63, 0x6b, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x10, 0x03, 0x2a, 0x27, 0x0a, 0x08, 0x43, 0x6d, 0x64, 0x4c, 0x6f, 0x67, 0x69, 0x63,
	0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x65, 0x66, 0x43, 0x6d, 0x64, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x10,
	0x00, 0x12, 0x0a, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x10, 0x01, 0x2a, 0x32, 0x0a,
	0x0a, 0x4b, 0x69, 0x63, 0x6b, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x44,
	0x65, 0x66, 0x4b, 0x69, 0x63, 0x6b, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x10, 0x00, 0x12, 0x11,
	0x0a, 0x0d, 0x4f, 0x6e, 0x6c, 0x79, 0x4f, 0x6e, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x10,
	0x01, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dispatch_enum_proto_rawDescOnce sync.Once
	file_dispatch_enum_proto_rawDescData = file_dispatch_enum_proto_rawDesc
)

func file_dispatch_enum_proto_rawDescGZIP() []byte {
	file_dispatch_enum_proto_rawDescOnce.Do(func() {
		file_dispatch_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_dispatch_enum_proto_rawDescData)
	})
	return file_dispatch_enum_proto_rawDescData
}

var file_dispatch_enum_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_dispatch_enum_proto_goTypes = []interface{}{
	(Grp)(0),        // 0: pb.Grp
	(CmdSys)(0),     // 1: pb.CmdSys
	(CmdLogic)(0),   // 2: pb.CmdLogic
	(KickReason)(0), // 3: pb.KickReason
}
var file_dispatch_enum_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_dispatch_enum_proto_init() }
func file_dispatch_enum_proto_init() {
	if File_dispatch_enum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dispatch_enum_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dispatch_enum_proto_goTypes,
		DependencyIndexes: file_dispatch_enum_proto_depIdxs,
		EnumInfos:         file_dispatch_enum_proto_enumTypes,
	}.Build()
	File_dispatch_enum_proto = out.File
	file_dispatch_enum_proto_rawDesc = nil
	file_dispatch_enum_proto_goTypes = nil
	file_dispatch_enum_proto_depIdxs = nil
}
