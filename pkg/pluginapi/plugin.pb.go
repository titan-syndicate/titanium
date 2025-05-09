// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: pkg/pluginapi/proto/plugin.proto

package pluginapi

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Empty is used for requests that don't need any parameters
type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_pkg_pluginapi_proto_plugin_proto_rawDescGZIP(), []int{0}
}

// NameResponse is the response from the Name RPC
type NameResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NameResponse) Reset() {
	*x = NameResponse{}
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NameResponse) ProtoMessage() {}

func (x *NameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NameResponse.ProtoReflect.Descriptor instead.
func (*NameResponse) Descriptor() ([]byte, []int) {
	return file_pkg_pluginapi_proto_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *NameResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// VersionResponse is the response from the Version RPC
type VersionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Version       string                 `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionResponse.ProtoReflect.Descriptor instead.
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return file_pkg_pluginapi_proto_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *VersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// ExecuteRequest is the request for the Execute RPC
type ExecuteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Args          []string               `protobuf:"bytes,1,rep,name=args,proto3" json:"args,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExecuteRequest) Reset() {
	*x = ExecuteRequest{}
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecuteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteRequest) ProtoMessage() {}

func (x *ExecuteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteRequest.ProtoReflect.Descriptor instead.
func (*ExecuteRequest) Descriptor() ([]byte, []int) {
	return file_pkg_pluginapi_proto_plugin_proto_rawDescGZIP(), []int{3}
}

func (x *ExecuteRequest) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

// ExecuteResponse is the response from the Execute RPC
type ExecuteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExecuteResponse) Reset() {
	*x = ExecuteResponse{}
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecuteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteResponse) ProtoMessage() {}

func (x *ExecuteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_pluginapi_proto_plugin_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteResponse.ProtoReflect.Descriptor instead.
func (*ExecuteResponse) Descriptor() ([]byte, []int) {
	return file_pkg_pluginapi_proto_plugin_proto_rawDescGZIP(), []int{4}
}

func (x *ExecuteResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_pkg_pluginapi_proto_plugin_proto protoreflect.FileDescriptor

const file_pkg_pluginapi_proto_plugin_proto_rawDesc = "" +
	"\n" +
	" pkg/pluginapi/proto/plugin.proto\x12\tpluginapi\"\a\n" +
	"\x05Empty\"\"\n" +
	"\fNameResponse\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"+\n" +
	"\x0fVersionResponse\x12\x18\n" +
	"\aversion\x18\x01 \x01(\tR\aversion\"$\n" +
	"\x0eExecuteRequest\x12\x12\n" +
	"\x04args\x18\x01 \x03(\tR\x04args\")\n" +
	"\x0fExecuteResponse\x12\x16\n" +
	"\x06result\x18\x01 \x01(\tR\x06result2\xbc\x01\n" +
	"\x06Plugin\x123\n" +
	"\x04Name\x12\x10.pluginapi.Empty\x1a\x17.pluginapi.NameResponse\"\x00\x129\n" +
	"\aVersion\x12\x10.pluginapi.Empty\x1a\x1a.pluginapi.VersionResponse\"\x00\x12B\n" +
	"\aExecute\x12\x19.pluginapi.ExecuteRequest\x1a\x1a.pluginapi.ExecuteResponse\"\x00B3Z1github.com/titan-syndicate/titanium/pkg/pluginapib\x06proto3"

var (
	file_pkg_pluginapi_proto_plugin_proto_rawDescOnce sync.Once
	file_pkg_pluginapi_proto_plugin_proto_rawDescData []byte
)

func file_pkg_pluginapi_proto_plugin_proto_rawDescGZIP() []byte {
	file_pkg_pluginapi_proto_plugin_proto_rawDescOnce.Do(func() {
		file_pkg_pluginapi_proto_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_pluginapi_proto_plugin_proto_rawDesc), len(file_pkg_pluginapi_proto_plugin_proto_rawDesc)))
	})
	return file_pkg_pluginapi_proto_plugin_proto_rawDescData
}

var file_pkg_pluginapi_proto_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_pluginapi_proto_plugin_proto_goTypes = []any{
	(*Empty)(nil),           // 0: pluginapi.Empty
	(*NameResponse)(nil),    // 1: pluginapi.NameResponse
	(*VersionResponse)(nil), // 2: pluginapi.VersionResponse
	(*ExecuteRequest)(nil),  // 3: pluginapi.ExecuteRequest
	(*ExecuteResponse)(nil), // 4: pluginapi.ExecuteResponse
}
var file_pkg_pluginapi_proto_plugin_proto_depIdxs = []int32{
	0, // 0: pluginapi.Plugin.Name:input_type -> pluginapi.Empty
	0, // 1: pluginapi.Plugin.Version:input_type -> pluginapi.Empty
	3, // 2: pluginapi.Plugin.Execute:input_type -> pluginapi.ExecuteRequest
	1, // 3: pluginapi.Plugin.Name:output_type -> pluginapi.NameResponse
	2, // 4: pluginapi.Plugin.Version:output_type -> pluginapi.VersionResponse
	4, // 5: pluginapi.Plugin.Execute:output_type -> pluginapi.ExecuteResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_pluginapi_proto_plugin_proto_init() }
func file_pkg_pluginapi_proto_plugin_proto_init() {
	if File_pkg_pluginapi_proto_plugin_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_pluginapi_proto_plugin_proto_rawDesc), len(file_pkg_pluginapi_proto_plugin_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_pluginapi_proto_plugin_proto_goTypes,
		DependencyIndexes: file_pkg_pluginapi_proto_plugin_proto_depIdxs,
		MessageInfos:      file_pkg_pluginapi_proto_plugin_proto_msgTypes,
	}.Build()
	File_pkg_pluginapi_proto_plugin_proto = out.File
	file_pkg_pluginapi_proto_plugin_proto_goTypes = nil
	file_pkg_pluginapi_proto_plugin_proto_depIdxs = nil
}
