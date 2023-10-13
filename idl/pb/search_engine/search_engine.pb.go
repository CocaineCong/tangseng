// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: search_engine.proto

package search_engine

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

type SearchEngineRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @inject_tag:form:"query" uri:"query"
	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty" form:"query" uri:"query"`
}

func (x *SearchEngineRequest) Reset() {
	*x = SearchEngineRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_engine_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchEngineRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchEngineRequest) ProtoMessage() {}

func (x *SearchEngineRequest) ProtoReflect() protoreflect.Message {
	mi := &file_search_engine_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchEngineRequest.ProtoReflect.Descriptor instead.
func (*SearchEngineRequest) Descriptor() ([]byte, []int) {
	return file_search_engine_proto_rawDescGZIP(), []int{0}
}

func (x *SearchEngineRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type PostData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *PostData) Reset() {
	*x = PostData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_engine_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostData) ProtoMessage() {}

func (x *PostData) ProtoReflect() protoreflect.Message {
	mi := &file_search_engine_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostData.ProtoReflect.Descriptor instead.
func (*PostData) Descriptor() ([]byte, []int) {
	return file_search_engine_proto_rawDescGZIP(), []int{1}
}

func (x *PostData) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *PostData) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type SearchEngineList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlId int64   `protobuf:"varint,1,opt,name=url_id,json=urlId,proto3" json:"url_id,omitempty"`
	Desc  string  `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc,omitempty"`
	Url   string  `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Score float32 `protobuf:"fixed32,4,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *SearchEngineList) Reset() {
	*x = SearchEngineList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_engine_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchEngineList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchEngineList) ProtoMessage() {}

func (x *SearchEngineList) ProtoReflect() protoreflect.Message {
	mi := &file_search_engine_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchEngineList.ProtoReflect.Descriptor instead.
func (*SearchEngineList) Descriptor() ([]byte, []int) {
	return file_search_engine_proto_rawDescGZIP(), []int{2}
}

func (x *SearchEngineList) GetUrlId() int64 {
	if x != nil {
		return x.UrlId
	}
	return 0
}

func (x *SearchEngineList) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *SearchEngineList) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *SearchEngineList) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

type SearchEngineResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code                 int64               `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string              `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Count                int64               `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
	SearchEngineInfoList []*SearchEngineList `protobuf:"bytes,4,rep,name=search_engine_info_list,json=searchEngineInfoList,proto3" json:"search_engine_info_list,omitempty"`
	Data                 []string            `protobuf:"bytes,5,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *SearchEngineResponse) Reset() {
	*x = SearchEngineResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_engine_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchEngineResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchEngineResponse) ProtoMessage() {}

func (x *SearchEngineResponse) ProtoReflect() protoreflect.Message {
	mi := &file_search_engine_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchEngineResponse.ProtoReflect.Descriptor instead.
func (*SearchEngineResponse) Descriptor() ([]byte, []int) {
	return file_search_engine_proto_rawDescGZIP(), []int{3}
}

func (x *SearchEngineResponse) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *SearchEngineResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *SearchEngineResponse) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *SearchEngineResponse) GetSearchEngineInfoList() []*SearchEngineList {
	if x != nil {
		return x.SearchEngineInfoList
	}
	return nil
}

func (x *SearchEngineResponse) GetData() []string {
	if x != nil {
		return x.Data
	}
	return nil
}

type WordAssociationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code                int64    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                 string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	WordAssociationList []string `protobuf:"bytes,3,rep,name=word_association_list,json=wordAssociationList,proto3" json:"word_association_list,omitempty"`
	Data                string   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *WordAssociationResponse) Reset() {
	*x = WordAssociationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_engine_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WordAssociationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WordAssociationResponse) ProtoMessage() {}

func (x *WordAssociationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_search_engine_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WordAssociationResponse.ProtoReflect.Descriptor instead.
func (*WordAssociationResponse) Descriptor() ([]byte, []int) {
	return file_search_engine_proto_rawDescGZIP(), []int{4}
}

func (x *WordAssociationResponse) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *WordAssociationResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *WordAssociationResponse) GetWordAssociationList() []string {
	if x != nil {
		return x.WordAssociationList
	}
	return nil
}

func (x *WordAssociationResponse) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_search_engine_proto protoreflect.FileDescriptor

var file_search_engine_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a, 0x13, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x22, 0x32, 0x0a, 0x08, 0x50, 0x6f, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x65, 0x0a, 0x10, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x75, 0x72,
	0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x75, 0x72, 0x6c, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0xb0, 0x01,
	0x0a, 0x14, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x48, 0x0a, 0x17, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x14, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x87, 0x01, 0x0a, 0x17, 0x57, 0x6f, 0x72, 0x64, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d,
	0x73, 0x67, 0x12, 0x32, 0x0a, 0x15, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x61, 0x73, 0x73, 0x6f, 0x63,
	0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x13, 0x77, 0x6f, 0x72, 0x64, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x9b, 0x01, 0x0a, 0x13, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x41, 0x0a, 0x12, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x14, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0f, 0x57, 0x6f, 0x72, 0x64, 0x41, 0x73, 0x73,
	0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18,
	0x2e, 0x57, 0x6f, 0x72, 0x64, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x2f, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x3b, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_search_engine_proto_rawDescOnce sync.Once
	file_search_engine_proto_rawDescData = file_search_engine_proto_rawDesc
)

func file_search_engine_proto_rawDescGZIP() []byte {
	file_search_engine_proto_rawDescOnce.Do(func() {
		file_search_engine_proto_rawDescData = protoimpl.X.CompressGZIP(file_search_engine_proto_rawDescData)
	})
	return file_search_engine_proto_rawDescData
}

var file_search_engine_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_search_engine_proto_goTypes = []interface{}{
	(*SearchEngineRequest)(nil),     // 0: SearchEngineRequest
	(*PostData)(nil),                // 1: PostData
	(*SearchEngineList)(nil),        // 2: SearchEngineList
	(*SearchEngineResponse)(nil),    // 3: SearchEngineResponse
	(*WordAssociationResponse)(nil), // 4: WordAssociationResponse
}
var file_search_engine_proto_depIdxs = []int32{
	2, // 0: SearchEngineResponse.search_engine_info_list:type_name -> SearchEngineList
	0, // 1: SearchEngineService.SearchEngineSearch:input_type -> SearchEngineRequest
	0, // 2: SearchEngineService.WordAssociation:input_type -> SearchEngineRequest
	3, // 3: SearchEngineService.SearchEngineSearch:output_type -> SearchEngineResponse
	4, // 4: SearchEngineService.WordAssociation:output_type -> WordAssociationResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_search_engine_proto_init() }
func file_search_engine_proto_init() {
	if File_search_engine_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_search_engine_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchEngineRequest); i {
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
		file_search_engine_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostData); i {
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
		file_search_engine_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchEngineList); i {
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
		file_search_engine_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchEngineResponse); i {
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
		file_search_engine_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WordAssociationResponse); i {
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
			RawDescriptor: file_search_engine_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_search_engine_proto_goTypes,
		DependencyIndexes: file_search_engine_proto_depIdxs,
		MessageInfos:      file_search_engine_proto_msgTypes,
	}.Build()
	File_search_engine_proto = out.File
	file_search_engine_proto_rawDesc = nil
	file_search_engine_proto_goTypes = nil
	file_search_engine_proto_depIdxs = nil
}
