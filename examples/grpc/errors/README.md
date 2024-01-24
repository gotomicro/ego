# protoc-gen-go-errors 使用说明

## 使用准备
- 需要本地提前安装好 protoc google/protobuf 等
- 需要下载最新的 protoc-gen-go-errors 插件，并添加到本地环境变量 

## 定义错误 proto 文件
```protobuf
syntax = "proto3";

package biz.v1;

// 下行注解是必要的，有这个注解 protoc-gen-go-errors 才尝试解析当前 protobuf 文件中的 enum，并基于 enum 生成错误桩代码
// @plugins=protoc-gen-go-errors

// language-specified package name
option go_package = "biz/v1;bizv1";
option java_multiple_files = true;
option java_outer_classname = "BizProto";
option java_package = "com.ego.biz.v1";

// enum Err 定义了错误的不同枚举值，protoc-gen-go-errors 会基于 enum Err 枚举值生成错误桩代码
// @code 为错误关联的gRPC Code (遵循 https://grpc.github.io/grpc/core/md_doc_statuscodes.html 定义，需要全大写)，
//       包含 OK、UNKNOWN、INVALID_ARGUMENT、PERMISSION_DENIED等
// @i18n.cn 国际化中文文案
// @i18n.en 国际化英文文案
enum Err {
  // 请求正常，实际上不算是一个错误
  // @code=OK
  // @i18n.cn="请求成功"
  // @i18n.en="OK"
  ERR_OK = 0;
  // 未知错误，比如业务panic了
  // @code=UNKNOWN             # 定义了这个错误关联的gRPC Code为：UNKNOWN
  // @i18n.cn="服务内部未知错误" # 定义了一个中文错误文案
  // @i18n.en="unknown error"  # 定义了一个英文错误文案
  ERR_UNKNOWN = 1;
  // 找不到指定用户
  // @code=NOT_FOUND
  // @i18n.cn="找不到指定用户"
  // @i18n.en="user not found"
  ERR_USER_NOT_FOUND = 2;
  // 用户ID不合法
  // @code=INVALID_ARGUMENT
  // @i18n.cn="用户ID不合法"
  // @i18n.en="invalid user id"
  ERR_USER_ID_NOT_VALID = 3;
}
```

## 使用
```bash
protoc --proto_path=. --go_out=paths=source_relative:. --go-errors_out=paths=source_relative:. ./errors.proto
```

详细用例可以参考 [./error.sh](./error.sh) 脚本