// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
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
	0x6f, 0x74, 0x6f, 0x12, 0x16, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x62, 0x69, 0x7a, 0x2f, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x15, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xa7, 0x08, 0x0a, 0x0c, 0x41, 0x70, 0x70, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x64, 0x0a, 0x09, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x41, 0x70, 0x70, 0x12, 0x29, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2c, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x41,
	0x6e, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x54, 0x0a,
	0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x80, 0x01, 0x0a, 0x15, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x42,
	0x61, 0x73, 0x69, 0x63, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12, 0x30, 0x2e,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x42, 0x61,
	0x73, 0x69, 0x63, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a,
	0x35, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c,
	0x42, 0x61, 0x73, 0x69, 0x63, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5e, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70,
	0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x12, 0x13, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x1a, 0x2f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41,
	0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x26, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x41, 0x70, 0x70, 0x12, 0x0c, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70,
	0x70, 0x1a, 0x0b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x4c,
	0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x2b, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a,
	0x0b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x72, 0x0a, 0x14,
	0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x41, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2c, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x47, 0x65,
	0x74, 0x41, 0x70, 0x70, 0x41, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x1a, 0x2c, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x70, 0x70, 0x41, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x48, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x25,
	0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e,
	0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x18, 0x52, 0x65,
	0x6c, 0x6f, 0x61, 0x64, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1b, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x1a, 0x0b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x73, 0x67,
	0x12, 0x34, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x52, 0x65, 0x6c,
	0x65, 0x61, 0x73, 0x65, 0x12, 0x13, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41,
	0x70, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x1a, 0x0b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x30, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x41, 0x70, 0x70,
	0x52, 0x65, 0x70, 0x6f, 0x12, 0x10, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41,
	0x70, 0x70, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x10, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x43, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x41,
	0x70, 0x70, 0x73, 0x42, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x10, 0x2e, 0x62, 0x69, 0x7a, 0x2e,
	0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x70, 0x6f, 0x1a, 0x20, 0x2e, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x51, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x42, 0x79, 0x52,
	0x65, 0x70, 0x6f, 0x12, 0x2d, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e,
	0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x70, 0x70, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x42, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x52,
	0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70,
	0x42, 0x0a, 0x5a, 0x08, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x70, 0x3b, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var file_api_app_app_proto_goTypes = []any{
	(*FileUploadRequest)(nil),             // 0: clusterruntime.api.app.FileUploadRequest
	(*emptypb.Empty)(nil),                 // 1: google.protobuf.Empty
	(*InstallBasicComponentReq)(nil),      // 2: clusterruntime.api.app.InstallBasicComponentReq
	(*biz.AppRelease)(nil),                // 3: biz.app.AppRelease
	(*biz.App)(nil),                       // 4: biz.app.App
	(*DeleteAppVersionReq)(nil),           // 5: clusterruntime.api.app.DeleteAppVersionReq
	(*GetAppAndVersionInfo)(nil),          // 6: clusterruntime.api.app.GetAppAndVersionInfo
	(*AppReleaseReq)(nil),                 // 7: clusterruntime.api.app.AppReleaseReq
	(*biz.AppReleaseResource)(nil),        // 8: biz.app.AppReleaseResource
	(*biz.AppRepo)(nil),                   // 9: biz.app.AppRepo
	(*GetAppDetailByRepoReq)(nil),         // 10: clusterruntime.api.app.GetAppDetailByRepoReq
	(*CheckClusterResponse)(nil),          // 11: clusterruntime.api.app.CheckClusterResponse
	(*InstallBasicComponentResponse)(nil), // 12: clusterruntime.api.app.InstallBasicComponentResponse
	(*AppReleaseResourceItems)(nil),       // 13: clusterruntime.api.app.AppReleaseResourceItems
	(*common.Msg)(nil),                    // 14: common.Msg
	(*AppItems)(nil),                      // 15: clusterruntime.api.app.AppItems
}
var file_api_app_app_proto_depIdxs = []int32{
	0,  // 0: clusterruntime.api.app.AppInterface.UploadApp:input_type -> clusterruntime.api.app.FileUploadRequest
	1,  // 1: clusterruntime.api.app.AppInterface.CheckCluster:input_type -> google.protobuf.Empty
	2,  // 2: clusterruntime.api.app.AppInterface.InstallBasicComponent:input_type -> clusterruntime.api.app.InstallBasicComponentReq
	3,  // 3: clusterruntime.api.app.AppInterface.GetAppReleaseResources:input_type -> biz.app.AppRelease
	4,  // 4: clusterruntime.api.app.AppInterface.DeleteApp:input_type -> biz.app.App
	5,  // 5: clusterruntime.api.app.AppInterface.DeleteAppVersion:input_type -> clusterruntime.api.app.DeleteAppVersionReq
	6,  // 6: clusterruntime.api.app.AppInterface.GetAppAndVersionInfo:input_type -> clusterruntime.api.app.GetAppAndVersionInfo
	7,  // 7: clusterruntime.api.app.AppInterface.AppRelease:input_type -> clusterruntime.api.app.AppReleaseReq
	8,  // 8: clusterruntime.api.app.AppInterface.ReloadAppReleaseResource:input_type -> biz.app.AppReleaseResource
	3,  // 9: clusterruntime.api.app.AppInterface.DeleteAppRelease:input_type -> biz.app.AppRelease
	9,  // 10: clusterruntime.api.app.AppInterface.AddAppRepo:input_type -> biz.app.AppRepo
	9,  // 11: clusterruntime.api.app.AppInterface.GetAppsByRepo:input_type -> biz.app.AppRepo
	10, // 12: clusterruntime.api.app.AppInterface.GetAppDetailByRepo:input_type -> clusterruntime.api.app.GetAppDetailByRepoReq
	6,  // 13: clusterruntime.api.app.AppInterface.UploadApp:output_type -> clusterruntime.api.app.GetAppAndVersionInfo
	11, // 14: clusterruntime.api.app.AppInterface.CheckCluster:output_type -> clusterruntime.api.app.CheckClusterResponse
	12, // 15: clusterruntime.api.app.AppInterface.InstallBasicComponent:output_type -> clusterruntime.api.app.InstallBasicComponentResponse
	13, // 16: clusterruntime.api.app.AppInterface.GetAppReleaseResources:output_type -> clusterruntime.api.app.AppReleaseResourceItems
	14, // 17: clusterruntime.api.app.AppInterface.DeleteApp:output_type -> common.Msg
	14, // 18: clusterruntime.api.app.AppInterface.DeleteAppVersion:output_type -> common.Msg
	6,  // 19: clusterruntime.api.app.AppInterface.GetAppAndVersionInfo:output_type -> clusterruntime.api.app.GetAppAndVersionInfo
	3,  // 20: clusterruntime.api.app.AppInterface.AppRelease:output_type -> biz.app.AppRelease
	14, // 21: clusterruntime.api.app.AppInterface.ReloadAppReleaseResource:output_type -> common.Msg
	14, // 22: clusterruntime.api.app.AppInterface.DeleteAppRelease:output_type -> common.Msg
	9,  // 23: clusterruntime.api.app.AppInterface.AddAppRepo:output_type -> biz.app.AppRepo
	15, // 24: clusterruntime.api.app.AppInterface.GetAppsByRepo:output_type -> clusterruntime.api.app.AppItems
	4,  // 25: clusterruntime.api.app.AppInterface.GetAppDetailByRepo:output_type -> biz.app.App
	13, // [13:26] is the sub-list for method output_type
	0,  // [0:13] is the sub-list for method input_type
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
