// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: api/app/app.proto

package app

import (
	common "github.com/f-rambo/cloud-copilot/cluster-runtime/api/common"
	biz "github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_api_app_app_proto protoreflect.FileDescriptor

var file_api_app_app_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x70, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x62, 0x69, 0x7a, 0x2f, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x32, 0xa5, 0x05, 0x0a, 0x0c, 0x41, 0x70, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x66, 0x61, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x70,
	0x70, 0x12, 0x16, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x70, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x41, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x41, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x61,
	0x70, 0x70, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x13,
	0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65,
	0x61, 0x73, 0x65, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c,
	0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x74, 0x65, 0x6d,
	0x73, 0x12, 0x26, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x12, 0x0c,
	0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x1a, 0x0b, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x39, 0x0a, 0x10, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x2e,
	0x61, 0x70, 0x70, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x0b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x4d, 0x73, 0x67, 0x12, 0x4c, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x41, 0x6e,
	0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x2e, 0x61,
	0x70, 0x70, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x41, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x47, 0x65,
	0x74, 0x41, 0x70, 0x70, 0x41, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x35, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x12, 0x12, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41,
	0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x10, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x13, 0x2e,
	0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x1a, 0x13, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70,
	0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x41, 0x70,
	0x70, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x10, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e,
	0x41, 0x70, 0x70, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x10, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70,
	0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x30, 0x0a, 0x0d, 0x47, 0x65, 0x74,
	0x41, 0x70, 0x70, 0x73, 0x42, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x10, 0x2e, 0x62, 0x69, 0x7a,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x0d, 0x2e, 0x61,
	0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x3e, 0x0a, 0x12, 0x47,
	0x65, 0x74, 0x41, 0x70, 0x70, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x42, 0x79, 0x52, 0x65, 0x70,
	0x6f, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x42, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e,
	0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x42, 0x0a, 0x5a, 0x08, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x70, 0x70, 0x3b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_app_app_proto_goTypes = []any{
	(*FileUploadRequest)(nil),       // 0: app.FileUploadRequest
	(*emptypb.Empty)(nil),           // 1: google.protobuf.Empty
	(*biz.AppRelease)(nil),          // 2: biz.app.AppRelease
	(*biz.App)(nil),                 // 3: biz.app.App
	(*DeleteAppVersionReq)(nil),     // 4: app.DeleteAppVersionReq
	(*GetAppAndVersionInfo)(nil),    // 5: app.GetAppAndVersionInfo
	(*AppReleaseReq)(nil),           // 6: app.AppReleaseReq
	(*biz.AppRepo)(nil),             // 7: biz.app.AppRepo
	(*GetAppDetailByRepoReq)(nil),   // 8: app.GetAppDetailByRepoReq
	(*CheckClusterResponse)(nil),    // 9: app.CheckClusterResponse
	(*AppReleaseResourceItems)(nil), // 10: app.AppReleaseResourceItems
	(*common.Msg)(nil),              // 11: common.Msg
	(*AppItems)(nil),                // 12: app.AppItems
}
var file_api_app_app_proto_depIdxs = []int32{
	0,  // 0: app.AppInterface.UploadApp:input_type -> app.FileUploadRequest
	1,  // 1: app.AppInterface.CheckCluster:input_type -> google.protobuf.Empty
	2,  // 2: app.AppInterface.GetClusterResources:input_type -> biz.app.AppRelease
	3,  // 3: app.AppInterface.DeleteApp:input_type -> biz.app.App
	4,  // 4: app.AppInterface.DeleteAppVersion:input_type -> app.DeleteAppVersionReq
	5,  // 5: app.AppInterface.GetAppAndVersionInfo:input_type -> app.GetAppAndVersionInfo
	6,  // 6: app.AppInterface.AppRelease:input_type -> app.AppReleaseReq
	2,  // 7: app.AppInterface.DeleteAppRelease:input_type -> biz.app.AppRelease
	7,  // 8: app.AppInterface.AddAppRepo:input_type -> biz.app.AppRepo
	7,  // 9: app.AppInterface.GetAppsByRepo:input_type -> biz.app.AppRepo
	8,  // 10: app.AppInterface.GetAppDetailByRepo:input_type -> app.GetAppDetailByRepoReq
	5,  // 11: app.AppInterface.UploadApp:output_type -> app.GetAppAndVersionInfo
	9,  // 12: app.AppInterface.CheckCluster:output_type -> app.CheckClusterResponse
	10, // 13: app.AppInterface.GetClusterResources:output_type -> app.AppReleaseResourceItems
	11, // 14: app.AppInterface.DeleteApp:output_type -> common.Msg
	11, // 15: app.AppInterface.DeleteAppVersion:output_type -> common.Msg
	5,  // 16: app.AppInterface.GetAppAndVersionInfo:output_type -> app.GetAppAndVersionInfo
	2,  // 17: app.AppInterface.AppRelease:output_type -> biz.app.AppRelease
	2,  // 18: app.AppInterface.DeleteAppRelease:output_type -> biz.app.AppRelease
	7,  // 19: app.AppInterface.AddAppRepo:output_type -> biz.app.AppRepo
	12, // 20: app.AppInterface.GetAppsByRepo:output_type -> app.AppItems
	3,  // 21: app.AppInterface.GetAppDetailByRepo:output_type -> biz.app.App
	11, // [11:22] is the sub-list for method output_type
	0,  // [0:11] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_app_app_proto_init() }
func file_api_app_app_proto_init() {
	if File_api_app_app_proto != nil {
		return
	}
	file_api_app_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_app_app_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_app_app_proto_goTypes,
		DependencyIndexes: file_api_app_app_proto_depIdxs,
	}.Build()
	File_api_app_app_proto = out.File
	file_api_app_app_proto_rawDesc = nil
	file_api_app_app_proto_goTypes = nil
	file_api_app_app_proto_depIdxs = nil
}