// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: errors.proto

package bizv1

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

// 错误定义
type Err int32

const (
	// 未知类型
	// @code=UNKNOWN
	Err_ERR_INVALID Err = 0
	// 找不到资源
	// @code=NOT_FOUND
	Err_ERR_USER_NOT_FOUND Err = 1
	// 客户端参数错误
	// @code=INVALID_ARGUMENT
	Err_ERR_CONTENT_MISSING Err = 2
)

// Enum value maps for Err.
var (
	Err_name = map[int32]string{
		0: "ERR_INVALID",
		1: "ERR_USER_NOT_FOUND",
		2: "ERR_CONTENT_MISSING",
	}
	Err_value = map[string]int32{
		"ERR_INVALID":         0,
		"ERR_USER_NOT_FOUND":  1,
		"ERR_CONTENT_MISSING": 2,
	}
)

func (x Err) Enum() *Err {
	p := new(Err)
	*p = x
	return p
}

func (x Err) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Err) Descriptor() protoreflect.EnumDescriptor {
	return file_errors_proto_enumTypes[0].Descriptor()
}

func (Err) Type() protoreflect.EnumType {
	return &file_errors_proto_enumTypes[0]
}

func (x Err) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Err.Descriptor instead.
func (Err) EnumDescriptor() ([]byte, []int) {
	return file_errors_proto_rawDescGZIP(), []int{0}
}

var File_errors_proto protoreflect.FileDescriptor

var file_errors_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x62, 0x69, 0x7a, 0x2e, 0x76, 0x31, 0x2a, 0x47, 0x0a, 0x03, 0x45, 0x72, 0x72, 0x12, 0x0f, 0x0a,
	0x0b, 0x45, 0x52, 0x52, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x16,
	0x0a, 0x12, 0x45, 0x52, 0x52, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46,
	0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x52, 0x52, 0x5f, 0x43, 0x4f,
	0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x42,
	0x2a, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x67, 0x6f, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x76,
	0x31, 0x42, 0x08, 0x42, 0x69, 0x7a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0c, 0x62,
	0x69, 0x7a, 0x2f, 0x76, 0x31, 0x3b, 0x62, 0x69, 0x7a, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_errors_proto_rawDescOnce sync.Once
	file_errors_proto_rawDescData = file_errors_proto_rawDesc
)

func file_errors_proto_rawDescGZIP() []byte {
	file_errors_proto_rawDescOnce.Do(func() {
		file_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_errors_proto_rawDescData)
	})
	return file_errors_proto_rawDescData
}

var file_errors_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errors_proto_goTypes = []interface{}{
	(Err)(0), // 0: biz.v1.Err
}
var file_errors_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_errors_proto_init() }
func file_errors_proto_init() {
	if File_errors_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_errors_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errors_proto_goTypes,
		DependencyIndexes: file_errors_proto_depIdxs,
		EnumInfos:         file_errors_proto_enumTypes,
	}.Build()
	File_errors_proto = out.File
	file_errors_proto_rawDesc = nil
	file_errors_proto_goTypes = nil
	file_errors_proto_depIdxs = nil
}