// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.20.2
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

// enum Err 定义了错误的不同枚举值，protoc-gen-go-errors 会基于 enum Err 枚举值生成错误桩代码
// @code 为错误关联的gRPC Code (遵循 https://grpc.github.io/grpc/core/md_doc_statuscodes.html 定义，需要全大写)，
//
//	包含 OK、UNKNOWN、INVALID_ARGUMENT、PERMISSION_DENIED等
//
// @i18n.cn 国际化中文文案
// @i18n.en 国际化英文文案
type Err int32

const (
	// 请求正常，实际上不算是一个错误
	// @code=OK
	// @i18n.cn="请求成功"
	// @i18n.en="OK"
	Err_ERR_OK Err = 0
	// 未知错误，比如业务panic了
	// @code=UNKNOWN             # 定义了这个错误关联的gRPC Code为：UNKNOWN
	// @i18n.cn="服务内部未知错误"        # 定义了一个中文错误文案
	// @i18n.en="unknown error"  # 定义了一个英文错误文案
	Err_ERR_UNKNOWN Err = 1
	// 找不到指定用户
	// @code=NOT_FOUND
	// @i18n.cn="找不到指定用户"
	// @i18n.en="user not found"
	Err_ERR_USER_NOT_FOUND Err = 2
	// 用户ID不合法
	// @code=INVALID_ARGUMENT
	// @i18n.cn="用户ID不合法"
	// @i18n.en="invalid user id"
	Err_ERR_USER_ID_NOT_VALID Err = 3
)

// Enum value maps for Err.
var (
	Err_name = map[int32]string{
		0: "ERR_OK",
		1: "ERR_UNKNOWN",
		2: "ERR_USER_NOT_FOUND",
		3: "ERR_USER_ID_NOT_VALID",
	}
	Err_value = map[string]int32{
		"ERR_OK":                0,
		"ERR_UNKNOWN":           1,
		"ERR_USER_NOT_FOUND":    2,
		"ERR_USER_ID_NOT_VALID": 3,
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
	0x62, 0x69, 0x7a, 0x2e, 0x76, 0x31, 0x2a, 0x55, 0x0a, 0x03, 0x45, 0x72, 0x72, 0x12, 0x0a, 0x0a,
	0x06, 0x45, 0x52, 0x52, 0x5f, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x52, 0x52,
	0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x52,
	0x52, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44,
	0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x45, 0x52, 0x52, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49,
	0x44, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x03, 0x42, 0x2a, 0x0a,
	0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x67, 0x6f, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x76, 0x31, 0x42,
	0x08, 0x42, 0x69, 0x7a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0c, 0x62, 0x69, 0x7a,
	0x2f, 0x76, 0x31, 0x3b, 0x62, 0x69, 0x7a, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
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
