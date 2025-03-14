// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: internal/biz/user.proto

package biz

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserStatus int32

const (
	UserStatus_USER_INIT    UserStatus = 0
	UserStatus_USER_ENABLE  UserStatus = 1
	UserStatus_USER_DISABLE UserStatus = 2
	UserStatus_USER_DELETED UserStatus = 3
)

// Enum value maps for UserStatus.
var (
	UserStatus_name = map[int32]string{
		0: "USER_INIT",
		1: "USER_ENABLE",
		2: "USER_DISABLE",
		3: "USER_DELETED",
	}
	UserStatus_value = map[string]int32{
		"USER_INIT":    0,
		"USER_ENABLE":  1,
		"USER_DISABLE": 2,
		"USER_DELETED": 3,
	}
)

func (x UserStatus) Enum() *UserStatus {
	p := new(UserStatus)
	*p = x
	return p
}

func (x UserStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_biz_user_proto_enumTypes[0].Descriptor()
}

func (UserStatus) Type() protoreflect.EnumType {
	return &file_internal_biz_user_proto_enumTypes[0]
}

func (x UserStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserStatus.Descriptor instead.
func (UserStatus) EnumDescriptor() ([]byte, []int) {
	return file_internal_biz_user_proto_rawDescGZIP(), []int{0}
}

type UserSignType int32

const (
	UserSignType_CREDENTIALS UserSignType = 0
	UserSignType_GITHUB      UserSignType = 1
)

// Enum value maps for UserSignType.
var (
	UserSignType_name = map[int32]string{
		0: "CREDENTIALS",
		1: "GITHUB",
	}
	UserSignType_value = map[string]int32{
		"CREDENTIALS": 0,
		"GITHUB":      1,
	}
)

func (x UserSignType) Enum() *UserSignType {
	p := new(UserSignType)
	*p = x
	return p
}

func (x UserSignType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserSignType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_biz_user_proto_enumTypes[1].Descriptor()
}

func (UserSignType) Type() protoreflect.EnumType {
	return &file_internal_biz_user_proto_enumTypes[1]
}

func (x UserSignType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserSignType.Descriptor instead.
func (UserSignType) EnumDescriptor() ([]byte, []int) {
	return file_internal_biz_user_proto_rawDescGZIP(), []int{1}
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                                        // @gotags: gorm:"column:id;primaryKey;AUTO_INCREMENT"
	Name        string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                                                     // @gotags: gorm:"column:name; default:''; NOT NULL"
	Email       string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`                                                   // @gotags: gorm:"column:email; default:''; NOT NULL"
	Password    string                 `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`                                             // @gotags: gorm:"column:password; default:''; NOT NULL"
	Status      UserStatus             `protobuf:"varint,5,opt,name=status,proto3,enum=biz.user.UserStatus" json:"status,omitempty"`                       // @gotags: gorm:"column:status; default:''; NOT NULL"
	AccessToken string                 `protobuf:"bytes,6,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`                    // @gotags: gorm:"-"`
	SignType    UserSignType           `protobuf:"varint,7,opt,name=sign_type,json=signType,proto3,enum=biz.user.UserSignType" json:"sign_type,omitempty"` // @gotags: gorm:"column:sign_type; default:''; NOT NULL"
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`                          // @gotags: gorm:"column:created_at; default:CURRENT_TIMESTAMP"
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`                          // @gotags: gorm:"column:updated_at; default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
	DeletedAt   *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`                         // @gotags: gorm:"column:deleted_at; default:NULL"
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_biz_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_internal_biz_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_internal_biz_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *User) GetStatus() UserStatus {
	if x != nil {
		return x.Status
	}
	return UserStatus_USER_INIT
}

func (x *User) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *User) GetSignType() UserSignType {
	if x != nil {
		return x.SignType
	}
	return UserSignType_CREDENTIALS
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *User) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

var File_internal_biz_user_proto protoreflect.FileDescriptor

var file_internal_biz_user_proto_rawDesc = []byte{
	0x0a, 0x17, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x69, 0x7a, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x62, 0x69, 0x7a, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x03, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x33, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x08, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x50, 0x0a, 0x0a, 0x55, 0x73,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x53, 0x45, 0x52,
	0x5f, 0x49, 0x4e, 0x49, 0x54, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x53, 0x45, 0x52, 0x5f,
	0x45, 0x4e, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x53, 0x45, 0x52,
	0x5f, 0x44, 0x49, 0x53, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x53,
	0x45, 0x52, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x03, 0x2a, 0x2b, 0x0a, 0x0c,
	0x55, 0x73, 0x65, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b,
	0x43, 0x52, 0x45, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x41, 0x4c, 0x53, 0x10, 0x00, 0x12, 0x0a, 0x0a,
	0x06, 0x47, 0x49, 0x54, 0x48, 0x55, 0x42, 0x10, 0x01, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x2d, 0x72, 0x61, 0x6d, 0x62, 0x6f, 0x2f,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2d, 0x63, 0x6f, 0x70, 0x69, 0x6c, 0x6f, 0x74, 0x2f, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x69, 0x7a, 0x3b, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_internal_biz_user_proto_rawDescOnce sync.Once
	file_internal_biz_user_proto_rawDescData = file_internal_biz_user_proto_rawDesc
)

func file_internal_biz_user_proto_rawDescGZIP() []byte {
	file_internal_biz_user_proto_rawDescOnce.Do(func() {
		file_internal_biz_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_biz_user_proto_rawDescData)
	})
	return file_internal_biz_user_proto_rawDescData
}

var file_internal_biz_user_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_internal_biz_user_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_biz_user_proto_goTypes = []any{
	(UserStatus)(0),               // 0: biz.user.UserStatus
	(UserSignType)(0),             // 1: biz.user.UserSignType
	(*User)(nil),                  // 2: biz.user.User
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_internal_biz_user_proto_depIdxs = []int32{
	0, // 0: biz.user.User.status:type_name -> biz.user.UserStatus
	1, // 1: biz.user.User.sign_type:type_name -> biz.user.UserSignType
	3, // 2: biz.user.User.created_at:type_name -> google.protobuf.Timestamp
	3, // 3: biz.user.User.updated_at:type_name -> google.protobuf.Timestamp
	3, // 4: biz.user.User.deleted_at:type_name -> google.protobuf.Timestamp
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_internal_biz_user_proto_init() }
func file_internal_biz_user_proto_init() {
	if File_internal_biz_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_biz_user_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*User); i {
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
			RawDescriptor: file_internal_biz_user_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_biz_user_proto_goTypes,
		DependencyIndexes: file_internal_biz_user_proto_depIdxs,
		EnumInfos:         file_internal_biz_user_proto_enumTypes,
		MessageInfos:      file_internal_biz_user_proto_msgTypes,
	}.Build()
	File_internal_biz_user_proto = out.File
	file_internal_biz_user_proto_rawDesc = nil
	file_internal_biz_user_proto_goTypes = nil
	file_internal_biz_user_proto_depIdxs = nil
}
