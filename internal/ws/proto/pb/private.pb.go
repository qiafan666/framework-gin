// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.3
// source: private.proto

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

// grp:2 cmd:1 健康检查
type ReqHealth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty" bson:"msg"`
}

func (x *ReqHealth) Reset() {
	*x = ReqHealth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_private_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqHealth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqHealth) ProtoMessage() {}

func (x *ReqHealth) ProtoReflect() protoreflect.Message {
	mi := &file_private_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqHealth.ProtoReflect.Descriptor instead.
func (*ReqHealth) Descriptor() ([]byte, []int) {
	return file_private_proto_rawDescGZIP(), []int{0}
}

func (x *ReqHealth) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type RspHealth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty" bson:"msg"`
}

func (x *RspHealth) Reset() {
	*x = RspHealth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_private_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RspHealth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspHealth) ProtoMessage() {}

func (x *RspHealth) ProtoReflect() protoreflect.Message {
	mi := &file_private_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspHealth.ProtoReflect.Descriptor instead.
func (*RspHealth) Descriptor() ([]byte, []int) {
	return file_private_proto_rawDescGZIP(), []int{1}
}

func (x *RspHealth) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_private_proto protoreflect.FileDescriptor

var file_private_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x22, 0x1d, 0x0a, 0x09, 0x52, 0x65, 0x71, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d,
	0x73, 0x67, 0x22, 0x1d, 0x0a, 0x09, 0x52, 0x73, 0x70, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_private_proto_rawDescOnce sync.Once
	file_private_proto_rawDescData = file_private_proto_rawDesc
)

func file_private_proto_rawDescGZIP() []byte {
	file_private_proto_rawDescOnce.Do(func() {
		file_private_proto_rawDescData = protoimpl.X.CompressGZIP(file_private_proto_rawDescData)
	})
	return file_private_proto_rawDescData
}

var file_private_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_private_proto_goTypes = []interface{}{
	(*ReqHealth)(nil), // 0: pb.ReqHealth
	(*RspHealth)(nil), // 1: pb.RspHealth
}
var file_private_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_private_proto_init() }
func file_private_proto_init() {
	if File_private_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_private_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqHealth); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_private_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RspHealth); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_private_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_private_proto_goTypes,
		DependencyIndexes: file_private_proto_depIdxs,
		MessageInfos:      file_private_proto_msgTypes,
	}.Build()
	File_private_proto = out.File
	file_private_proto_rawDesc = nil
	file_private_proto_goTypes = nil
	file_private_proto_depIdxs = nil
}
