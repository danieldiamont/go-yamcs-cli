// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: yamcs/protobuf/packets/packets.proto

package packets

import (
	protobuf "github.com/danieldiamont/go-yamcs-cli/yamcs/protobuf"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type TmPacketData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Raw packet binary
	Packet         []byte                  `protobuf:"bytes,2,req,name=packet" json:"packet,omitempty"`
	SequenceNumber *int32                  `protobuf:"varint,4,opt,name=sequenceNumber" json:"sequenceNumber,omitempty"`
	Id             *protobuf.NamedObjectId `protobuf:"bytes,5,opt,name=id" json:"id,omitempty"`
	// When the packet was generated
	GenerationTime *timestamp.Timestamp `protobuf:"bytes,9,opt,name=generationTime" json:"generationTime,omitempty"`
	// When the signal has been received on the groun
	EarthReceptionTime *timestamp.Timestamp `protobuf:"bytes,10,opt,name=earthReceptionTime" json:"earthReceptionTime,omitempty"`
	// When the packet was received by Yamcs
	ReceptionTime *timestamp.Timestamp `protobuf:"bytes,8,opt,name=receptionTime" json:"receptionTime,omitempty"`
	// Name of the Yamcs link where this packet was received from
	Link *string `protobuf:"bytes,11,opt,name=link" json:"link,omitempty"`
}

func (x *TmPacketData) Reset() {
	*x = TmPacketData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_yamcs_protobuf_packets_packets_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TmPacketData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TmPacketData) ProtoMessage() {}

func (x *TmPacketData) ProtoReflect() protoreflect.Message {
	mi := &file_yamcs_protobuf_packets_packets_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TmPacketData.ProtoReflect.Descriptor instead.
func (*TmPacketData) Descriptor() ([]byte, []int) {
	return file_yamcs_protobuf_packets_packets_proto_rawDescGZIP(), []int{0}
}

func (x *TmPacketData) GetPacket() []byte {
	if x != nil {
		return x.Packet
	}
	return nil
}

func (x *TmPacketData) GetSequenceNumber() int32 {
	if x != nil && x.SequenceNumber != nil {
		return *x.SequenceNumber
	}
	return 0
}

func (x *TmPacketData) GetId() *protobuf.NamedObjectId {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *TmPacketData) GetGenerationTime() *timestamp.Timestamp {
	if x != nil {
		return x.GenerationTime
	}
	return nil
}

func (x *TmPacketData) GetEarthReceptionTime() *timestamp.Timestamp {
	if x != nil {
		return x.EarthReceptionTime
	}
	return nil
}

func (x *TmPacketData) GetReceptionTime() *timestamp.Timestamp {
	if x != nil {
		return x.ReceptionTime
	}
	return nil
}

func (x *TmPacketData) GetLink() string {
	if x != nil && x.Link != nil {
		return *x.Link
	}
	return ""
}

var File_yamcs_protobuf_packets_packets_proto protoreflect.FileDescriptor

var file_yamcs_protobuf_packets_packets_proto_rawDesc = []byte{
	0x0a, 0x24, 0x79, 0x61, 0x6d, 0x63, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x79, 0x61, 0x6d, 0x63, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1a, 0x79, 0x61, 0x6d, 0x63, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x79, 0x61, 0x6d, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe3, 0x02, 0x0a, 0x0c,
	0x54, 0x6d, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x79, 0x61, 0x6d, 0x63, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x64, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0e, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x4a, 0x0a, 0x12, 0x65, 0x61, 0x72, 0x74, 0x68, 0x52, 0x65, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x12, 0x65, 0x61, 0x72, 0x74, 0x68, 0x52, 0x65,
	0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x0d, 0x72,
	0x65, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0d,
	0x72, 0x65, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x69, 0x6e,
	0x6b, 0x42, 0x62, 0x0a, 0x12, 0x6f, 0x72, 0x67, 0x2e, 0x79, 0x61, 0x6d, 0x63, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x42, 0x0c, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x6e, 0x69, 0x65, 0x6c, 0x64, 0x69, 0x61, 0x6d, 0x6f, 0x6e,
	0x74, 0x2f, 0x67, 0x6f, 0x2d, 0x79, 0x61, 0x6d, 0x63, 0x73, 0x2d, 0x63, 0x6c, 0x69, 0x2f, 0x79,
	0x61, 0x6d, 0x63, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73,
}

var (
	file_yamcs_protobuf_packets_packets_proto_rawDescOnce sync.Once
	file_yamcs_protobuf_packets_packets_proto_rawDescData = file_yamcs_protobuf_packets_packets_proto_rawDesc
)

func file_yamcs_protobuf_packets_packets_proto_rawDescGZIP() []byte {
	file_yamcs_protobuf_packets_packets_proto_rawDescOnce.Do(func() {
		file_yamcs_protobuf_packets_packets_proto_rawDescData = protoimpl.X.CompressGZIP(file_yamcs_protobuf_packets_packets_proto_rawDescData)
	})
	return file_yamcs_protobuf_packets_packets_proto_rawDescData
}

var file_yamcs_protobuf_packets_packets_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_yamcs_protobuf_packets_packets_proto_goTypes = []interface{}{
	(*TmPacketData)(nil),           // 0: yamcs.protobuf.packets.TmPacketData
	(*protobuf.NamedObjectId)(nil), // 1: yamcs.protobuf.NamedObjectId
	(*timestamp.Timestamp)(nil),    // 2: google.protobuf.Timestamp
}
var file_yamcs_protobuf_packets_packets_proto_depIdxs = []int32{
	1, // 0: yamcs.protobuf.packets.TmPacketData.id:type_name -> yamcs.protobuf.NamedObjectId
	2, // 1: yamcs.protobuf.packets.TmPacketData.generationTime:type_name -> google.protobuf.Timestamp
	2, // 2: yamcs.protobuf.packets.TmPacketData.earthReceptionTime:type_name -> google.protobuf.Timestamp
	2, // 3: yamcs.protobuf.packets.TmPacketData.receptionTime:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_yamcs_protobuf_packets_packets_proto_init() }
func file_yamcs_protobuf_packets_packets_proto_init() {
	if File_yamcs_protobuf_packets_packets_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_yamcs_protobuf_packets_packets_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TmPacketData); i {
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
			RawDescriptor: file_yamcs_protobuf_packets_packets_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_yamcs_protobuf_packets_packets_proto_goTypes,
		DependencyIndexes: file_yamcs_protobuf_packets_packets_proto_depIdxs,
		MessageInfos:      file_yamcs_protobuf_packets_packets_proto_msgTypes,
	}.Build()
	File_yamcs_protobuf_packets_packets_proto = out.File
	file_yamcs_protobuf_packets_packets_proto_rawDesc = nil
	file_yamcs_protobuf_packets_packets_proto_goTypes = nil
	file_yamcs_protobuf_packets_packets_proto_depIdxs = nil
}
