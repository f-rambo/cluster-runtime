// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: api/cluster/cluster.proto

package cluster

import (
	biz "github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_api_cluster_cluster_proto protoreflect.FileDescriptor

var file_api_cluster_cluster_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x1a, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x62, 0x69, 0x7a, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xa5,
	0x02, 0x0a, 0x10, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x15, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x12, 0x14, 0x2e, 0x62,
	0x69, 0x7a, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x1a, 0x2c, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64,
	0x12, 0x3c, 0x0a, 0x0e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x14, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x1a, 0x14, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x3a,
	0x0a, 0x0c, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x14,
	0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x1a, 0x14, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x0c, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x62, 0x69, 0x7a,
	0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x1a, 0x14, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x42, 0x0e, 0x5a, 0x0c, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x3b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_cluster_cluster_proto_goTypes = []any{
	(*biz.Cluster)(nil),      // 0: biz.cluster.Cluster
	(*ClusterInstalled)(nil), // 1: clusterruntime.api.cluster.ClusterInstalled
}
var file_api_cluster_cluster_proto_depIdxs = []int32{
	0, // 0: clusterruntime.api.cluster.ClusterInterface.CheckClusterInstalled:input_type -> biz.cluster.Cluster
	0, // 1: clusterruntime.api.cluster.ClusterInterface.CurrentCluster:input_type -> biz.cluster.Cluster
	0, // 2: clusterruntime.api.cluster.ClusterInterface.HandlerNodes:input_type -> biz.cluster.Cluster
	0, // 3: clusterruntime.api.cluster.ClusterInterface.StartCluster:input_type -> biz.cluster.Cluster
	1, // 4: clusterruntime.api.cluster.ClusterInterface.CheckClusterInstalled:output_type -> clusterruntime.api.cluster.ClusterInstalled
	0, // 5: clusterruntime.api.cluster.ClusterInterface.CurrentCluster:output_type -> biz.cluster.Cluster
	0, // 6: clusterruntime.api.cluster.ClusterInterface.HandlerNodes:output_type -> biz.cluster.Cluster
	0, // 7: clusterruntime.api.cluster.ClusterInterface.StartCluster:output_type -> biz.cluster.Cluster
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_cluster_cluster_proto_init() }
func file_api_cluster_cluster_proto_init() {
	if File_api_cluster_cluster_proto != nil {
		return
	}
	file_api_cluster_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_cluster_cluster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_cluster_cluster_proto_goTypes,
		DependencyIndexes: file_api_cluster_cluster_proto_depIdxs,
	}.Build()
	File_api_cluster_cluster_proto = out.File
	file_api_cluster_cluster_proto_rawDesc = nil
	file_api_cluster_cluster_proto_goTypes = nil
	file_api_cluster_cluster_proto_depIdxs = nil
}
