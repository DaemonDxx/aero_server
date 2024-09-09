// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.1
// source: auth/auth.v1.proto

package authpb

import (
	user "github.com/daemondxx/lks_back/gen/pb/go/user"
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

type AuthSystem int32

const (
	AuthSystem_SYSTEM_UNKNOWN AuthSystem = 0
	AuthSystem_SYSTEM_ACCORD  AuthSystem = 1
	AuthSystem_SYSTEM_LKS     AuthSystem = 2
)

// Enum value maps for AuthSystem.
var (
	AuthSystem_name = map[int32]string{
		0: "SYSTEM_UNKNOWN",
		1: "SYSTEM_ACCORD",
		2: "SYSTEM_LKS",
	}
	AuthSystem_value = map[string]int32{
		"SYSTEM_UNKNOWN": 0,
		"SYSTEM_ACCORD":  1,
		"SYSTEM_LKS":     2,
	}
)

func (x AuthSystem) Enum() *AuthSystem {
	p := new(AuthSystem)
	*p = x
	return p
}

func (x AuthSystem) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuthSystem) Descriptor() protoreflect.EnumDescriptor {
	return file_auth_auth_v1_proto_enumTypes[0].Descriptor()
}

func (AuthSystem) Type() protoreflect.EnumType {
	return &file_auth_auth_v1_proto_enumTypes[0]
}

func (x AuthSystem) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuthSystem.Descriptor instead.
func (AuthSystem) EnumDescriptor() ([]byte, []int) {
	return file_auth_auth_v1_proto_rawDescGZIP(), []int{0}
}

type AuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *user.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *AuthResponse) Reset() {
	*x = AuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_v1_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResponse) ProtoMessage() {}

func (x *AuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_v1_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResponse.ProtoReflect.Descriptor instead.
func (*AuthResponse) Descriptor() ([]byte, []int) {
	return file_auth_auth_v1_proto_rawDescGZIP(), []int{0}
}

func (x *AuthResponse) GetUser() *user.User {
	if x != nil {
		return x.User
	}
	return nil
}

type ErrorDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	System AuthSystem `protobuf:"varint,1,opt,name=system,proto3,enum=auth.v1.AuthSystem" json:"system,omitempty"`
}

func (x *ErrorDetails) Reset() {
	*x = ErrorDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_v1_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorDetails) ProtoMessage() {}

func (x *ErrorDetails) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_v1_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorDetails.ProtoReflect.Descriptor instead.
func (*ErrorDetails) Descriptor() ([]byte, []int) {
	return file_auth_auth_v1_proto_rawDescGZIP(), []int{1}
}

func (x *ErrorDetails) GetSystem() AuthSystem {
	if x != nil {
		return x.System
	}
	return AuthSystem_SYSTEM_UNKNOWN
}

type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Details *ErrorDetails `protobuf:"bytes,1,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_v1_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_v1_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_auth_auth_v1_proto_rawDescGZIP(), []int{2}
}

func (x *CheckResponse) GetDetails() *ErrorDetails {
	if x != nil {
		return x.Details
	}
	return nil
}

var File_auth_auth_v1_proto protoreflect.FileDescriptor

var file_auth_auth_v1_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x1a, 0x12, 0x75,
	0x73, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x31, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x21, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x22, 0x3b, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x22, 0x40, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x2a, 0x43, 0x0a, 0x0a, 0x41, 0x75, 0x74, 0x68, 0x53, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f,
	0x41, 0x43, 0x43, 0x4f, 0x52, 0x44, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x59, 0x53, 0x54,
	0x45, 0x4d, 0x5f, 0x4c, 0x4b, 0x53, 0x10, 0x02, 0x32, 0x73, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12,
	0x11, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x1a, 0x15, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x12, 0x11, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x35, 0x5a,
	0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x65, 0x6d,
	0x6f, 0x6e, 0x64, 0x78, 0x78, 0x2f, 0x6c, 0x6b, 0x73, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x3b, 0x61, 0x75,
	0x74, 0x68, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_auth_v1_proto_rawDescOnce sync.Once
	file_auth_auth_v1_proto_rawDescData = file_auth_auth_v1_proto_rawDesc
)

func file_auth_auth_v1_proto_rawDescGZIP() []byte {
	file_auth_auth_v1_proto_rawDescOnce.Do(func() {
		file_auth_auth_v1_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_auth_v1_proto_rawDescData)
	})
	return file_auth_auth_v1_proto_rawDescData
}

var file_auth_auth_v1_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_auth_auth_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_auth_auth_v1_proto_goTypes = []interface{}{
	(AuthSystem)(0),       // 0: auth.v1.AuthSystem
	(*AuthResponse)(nil),  // 1: auth.v1.AuthResponse
	(*ErrorDetails)(nil),  // 2: auth.v1.ErrorDetails
	(*CheckResponse)(nil), // 3: auth.v1.CheckResponse
	(*user.User)(nil),     // 4: user.v1.User
	(*user.UserInfo)(nil), // 5: user.v1.UserInfo
}
var file_auth_auth_v1_proto_depIdxs = []int32{
	4, // 0: auth.v1.AuthResponse.user:type_name -> user.v1.User
	0, // 1: auth.v1.ErrorDetails.system:type_name -> auth.v1.AuthSystem
	2, // 2: auth.v1.CheckResponse.details:type_name -> auth.v1.ErrorDetails
	5, // 3: auth.v1.AuthService.Auth:input_type -> user.v1.UserInfo
	5, // 4: auth.v1.AuthService.Check:input_type -> user.v1.UserInfo
	1, // 5: auth.v1.AuthService.Auth:output_type -> auth.v1.AuthResponse
	3, // 6: auth.v1.AuthService.Check:output_type -> auth.v1.CheckResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_auth_auth_v1_proto_init() }
func file_auth_auth_v1_proto_init() {
	if File_auth_auth_v1_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_auth_v1_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthResponse); i {
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
		file_auth_auth_v1_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorDetails); i {
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
		file_auth_auth_v1_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
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
			RawDescriptor: file_auth_auth_v1_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_auth_v1_proto_goTypes,
		DependencyIndexes: file_auth_auth_v1_proto_depIdxs,
		EnumInfos:         file_auth_auth_v1_proto_enumTypes,
		MessageInfos:      file_auth_auth_v1_proto_msgTypes,
	}.Build()
	File_auth_auth_v1_proto = out.File
	file_auth_auth_v1_proto_rawDesc = nil
	file_auth_auth_v1_proto_goTypes = nil
	file_auth_auth_v1_proto_depIdxs = nil
}
