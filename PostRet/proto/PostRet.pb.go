// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: proto/PostRet.proto

package PostRet

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

type PostRetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mobile   string `protobuf:"bytes,1,opt,name=mobile,proto3" json:"mobile,omitempty"`
	SmsCode  string `protobuf:"bytes,2,opt,name=sms_code,json=smsCode,proto3" json:"sms_code,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *PostRetRequest) Reset() {
	*x = PostRetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_PostRet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostRetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostRetRequest) ProtoMessage() {}

func (x *PostRetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_PostRet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostRetRequest.ProtoReflect.Descriptor instead.
func (*PostRetRequest) Descriptor() ([]byte, []int) {
	return file_proto_PostRet_proto_rawDescGZIP(), []int{0}
}

func (x *PostRetRequest) GetMobile() string {
	if x != nil {
		return x.Mobile
	}
	return ""
}

func (x *PostRetRequest) GetSmsCode() string {
	if x != nil {
		return x.SmsCode
	}
	return ""
}

func (x *PostRetRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type PostRetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error     string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Errmsg    string `protobuf:"bytes,2,opt,name=errmsg,proto3" json:"errmsg,omitempty"`
	SessionId string `protobuf:"bytes,3,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *PostRetResponse) Reset() {
	*x = PostRetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_PostRet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostRetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostRetResponse) ProtoMessage() {}

func (x *PostRetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_PostRet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostRetResponse.ProtoReflect.Descriptor instead.
func (*PostRetResponse) Descriptor() ([]byte, []int) {
	return file_proto_PostRet_proto_rawDescGZIP(), []int{1}
}

func (x *PostRetResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *PostRetResponse) GetErrmsg() string {
	if x != nil {
		return x.Errmsg
	}
	return ""
}

func (x *PostRetResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

var File_proto_PostRet_proto protoreflect.FileDescriptor

var file_proto_PostRet_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x22, 0x5f,
	0x0a, 0x0e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x6d, 0x73, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x6d, 0x73, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0x5e, 0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6d,
	0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6d, 0x73, 0x67,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x32,
	0x49, 0x0a, 0x07, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x12, 0x3e, 0x0a, 0x07, 0x50, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x74, 0x12, 0x17, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x2e,
	0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18,
	0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_PostRet_proto_rawDescOnce sync.Once
	file_proto_PostRet_proto_rawDescData = file_proto_PostRet_proto_rawDesc
)

func file_proto_PostRet_proto_rawDescGZIP() []byte {
	file_proto_PostRet_proto_rawDescOnce.Do(func() {
		file_proto_PostRet_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_PostRet_proto_rawDescData)
	})
	return file_proto_PostRet_proto_rawDescData
}

var file_proto_PostRet_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_PostRet_proto_goTypes = []interface{}{
	(*PostRetRequest)(nil),  // 0: PostRet.PostRetRequest
	(*PostRetResponse)(nil), // 1: PostRet.PostRetResponse
}
var file_proto_PostRet_proto_depIdxs = []int32{
	0, // 0: PostRet.PostRet.PostRet:input_type -> PostRet.PostRetRequest
	1, // 1: PostRet.PostRet.PostRet:output_type -> PostRet.PostRetResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_PostRet_proto_init() }
func file_proto_PostRet_proto_init() {
	if File_proto_PostRet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_PostRet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostRetRequest); i {
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
		file_proto_PostRet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostRetResponse); i {
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
			RawDescriptor: file_proto_PostRet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_PostRet_proto_goTypes,
		DependencyIndexes: file_proto_PostRet_proto_depIdxs,
		MessageInfos:      file_proto_PostRet_proto_msgTypes,
	}.Build()
	File_proto_PostRet_proto = out.File
	file_proto_PostRet_proto_rawDesc = nil
	file_proto_PostRet_proto_goTypes = nil
	file_proto_PostRet_proto_depIdxs = nil
}
