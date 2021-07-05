// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.13.0
// source: locker.proto

package client

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

type Entity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// идентификатор владения
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Entity) Reset() {
	*x = Entity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_locker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity) ProtoMessage() {}

func (x *Entity) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity.ProtoReflect.Descriptor instead.
func (*Entity) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{0}
}

func (x *Entity) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_locker_proto protoreflect.FileDescriptor

var file_locker_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16,
	0x72, 0x75, 0x2e, 0x73, 0x6f, 0x66, 0x74, 0x6c, 0x79, 0x6e, 0x78, 0x2e, 0x6c, 0x6f, 0x63, 0x6b,
	0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x22, 0x1c, 0x0a, 0x06, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x32, 0x9f, 0x01, 0x0a, 0x06, 0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x12,
	0x48, 0x0a, 0x04, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x1e, 0x2e, 0x72, 0x75, 0x2e, 0x73, 0x6f, 0x66,
	0x74, 0x6c, 0x79, 0x6e, 0x78, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x1e, 0x2e, 0x72, 0x75, 0x2e, 0x73, 0x6f, 0x66,
	0x74, 0x6c, 0x79, 0x6e, 0x78, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x30, 0x01, 0x12, 0x4b, 0x0a, 0x07, 0x54, 0x72, 0x79,
	0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x1e, 0x2e, 0x72, 0x75, 0x2e, 0x73, 0x6f, 0x66, 0x74, 0x6c, 0x79,
	0x6e, 0x78, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x1a, 0x1e, 0x2e, 0x72, 0x75, 0x2e, 0x73, 0x6f, 0x66, 0x74, 0x6c, 0x79,
	0x6e, 0x78, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x30, 0x01, 0x42, 0x55, 0x0a, 0x16, 0x72, 0x75, 0x2e, 0x73, 0x6f, 0x66,
	0x74, 0x6c, 0x79, 0x6e, 0x78, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69,
	0x42, 0x0a, 0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x50, 0x00, 0x5a, 0x20,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x66, 0x74, 0x6b,
	0x6f, 0x74, 0x2f, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0xa2, 0x02, 0x0a, 0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_locker_proto_rawDescOnce sync.Once
	file_locker_proto_rawDescData = file_locker_proto_rawDesc
)

func file_locker_proto_rawDescGZIP() []byte {
	file_locker_proto_rawDescOnce.Do(func() {
		file_locker_proto_rawDescData = protoimpl.X.CompressGZIP(file_locker_proto_rawDescData)
	})
	return file_locker_proto_rawDescData
}

var file_locker_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_locker_proto_goTypes = []interface{}{
	(*Entity)(nil), // 0: ru.softlynx.locker.api.Entity
}
var file_locker_proto_depIdxs = []int32{
	0, // 0: ru.softlynx.locker.api.Locker.Lock:input_type -> ru.softlynx.locker.api.Entity
	0, // 1: ru.softlynx.locker.api.Locker.TryLock:input_type -> ru.softlynx.locker.api.Entity
	0, // 2: ru.softlynx.locker.api.Locker.Lock:output_type -> ru.softlynx.locker.api.Entity
	0, // 3: ru.softlynx.locker.api.Locker.TryLock:output_type -> ru.softlynx.locker.api.Entity
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_locker_proto_init() }
func file_locker_proto_init() {
	if File_locker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_locker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entity); i {
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
			RawDescriptor: file_locker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_locker_proto_goTypes,
		DependencyIndexes: file_locker_proto_depIdxs,
		MessageInfos:      file_locker_proto_msgTypes,
	}.Build()
	File_locker_proto = out.File
	file_locker_proto_rawDesc = nil
	file_locker_proto_goTypes = nil
	file_locker_proto_depIdxs = nil
}
