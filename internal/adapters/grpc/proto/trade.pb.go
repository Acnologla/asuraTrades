// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.1
// source: proto/trade.proto

package proto

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

type GetUserProfileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserProfileRequest) Reset() {
	*x = GetUserProfileRequest{}
	mi := &file_proto_trade_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserProfileRequest) ProtoMessage() {}

func (x *GetUserProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_trade_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserProfileRequest.ProtoReflect.Descriptor instead.
func (*GetUserProfileRequest) Descriptor() ([]byte, []int) {
	return file_proto_trade_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserProfileRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type TradeItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type          int32                  `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	UserId        uint64                 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TradeItem) Reset() {
	*x = TradeItem{}
	mi := &file_proto_trade_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TradeItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradeItem) ProtoMessage() {}

func (x *TradeItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_trade_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradeItem.ProtoReflect.Descriptor instead.
func (*TradeItem) Descriptor() ([]byte, []int) {
	return file_proto_trade_proto_rawDescGZIP(), []int{1}
}

func (x *TradeItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TradeItem) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *TradeItem) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type FinishTradeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AuthorId      uint64                 `protobuf:"varint,1,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	OtherId       uint64                 `protobuf:"varint,2,opt,name=other_id,json=otherId,proto3" json:"other_id,omitempty"`
	Items         []*TradeItem           `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FinishTradeRequest) Reset() {
	*x = FinishTradeRequest{}
	mi := &file_proto_trade_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FinishTradeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishTradeRequest) ProtoMessage() {}

func (x *FinishTradeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_trade_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishTradeRequest.ProtoReflect.Descriptor instead.
func (*FinishTradeRequest) Descriptor() ([]byte, []int) {
	return file_proto_trade_proto_rawDescGZIP(), []int{2}
}

func (x *FinishTradeRequest) GetAuthorId() uint64 {
	if x != nil {
		return x.AuthorId
	}
	return 0
}

func (x *FinishTradeRequest) GetOtherId() uint64 {
	if x != nil {
		return x.OtherId
	}
	return 0
}

func (x *FinishTradeRequest) GetItems() []*TradeItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type Rooster struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        uint64                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Origin        string                 `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	Type          int32                  `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Rooster) Reset() {
	*x = Rooster{}
	mi := &file_proto_trade_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Rooster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rooster) ProtoMessage() {}

func (x *Rooster) ProtoReflect() protoreflect.Message {
	mi := &file_proto_trade_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rooster.ProtoReflect.Descriptor instead.
func (*Rooster) Descriptor() ([]byte, []int) {
	return file_proto_trade_proto_rawDescGZIP(), []int{3}
}

func (x *Rooster) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Rooster) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Rooster) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *Rooster) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

type Item struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        uint64                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Quantity      int32                  `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	ItemId        int32                  `protobuf:"varint,4,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	Type          int32                  `protobuf:"varint,5,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Item) Reset() {
	*x = Item{}
	mi := &file_proto_trade_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_proto_trade_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_proto_trade_proto_rawDescGZIP(), []int{4}
}

func (x *Item) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Item) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Item) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Item) GetItemId() int32 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *Item) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

type GetUserProfileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Roosters      []*Rooster             `protobuf:"bytes,2,rep,name=roosters,proto3" json:"roosters,omitempty"`
	Items         []*Item                `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserProfileResponse) Reset() {
	*x = GetUserProfileResponse{}
	mi := &file_proto_trade_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserProfileResponse) ProtoMessage() {}

func (x *GetUserProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_trade_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserProfileResponse.ProtoReflect.Descriptor instead.
func (*GetUserProfileResponse) Descriptor() ([]byte, []int) {
	return file_proto_trade_proto_rawDescGZIP(), []int{5}
}

func (x *GetUserProfileResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetUserProfileResponse) GetRoosters() []*Rooster {
	if x != nil {
		return x.Roosters
	}
	return nil
}

func (x *GetUserProfileResponse) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

type FinishTradeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ok            bool                   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Error         *string                `protobuf:"bytes,2,opt,name=error,proto3,oneof" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FinishTradeResponse) Reset() {
	*x = FinishTradeResponse{}
	mi := &file_proto_trade_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FinishTradeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishTradeResponse) ProtoMessage() {}

func (x *FinishTradeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_trade_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishTradeResponse.ProtoReflect.Descriptor instead.
func (*FinishTradeResponse) Descriptor() ([]byte, []int) {
	return file_proto_trade_proto_rawDescGZIP(), []int{6}
}

func (x *FinishTradeResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *FinishTradeResponse) GetError() string {
	if x != nil && x.Error != nil {
		return *x.Error
	}
	return ""
}

var File_proto_trade_proto protoreflect.FileDescriptor

const file_proto_trade_proto_rawDesc = "" +
	"\n" +
	"\x11proto/trade.proto\x12\x05trade\"'\n" +
	"\x15GetUserProfileRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\"H\n" +
	"\tTradeItem\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04type\x18\x02 \x01(\x05R\x04type\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\x04R\x06userId\"t\n" +
	"\x12FinishTradeRequest\x12\x1b\n" +
	"\tauthor_id\x18\x01 \x01(\x04R\bauthorId\x12\x19\n" +
	"\bother_id\x18\x02 \x01(\x04R\aotherId\x12&\n" +
	"\x05items\x18\x03 \x03(\v2\x10.trade.TradeItemR\x05items\"^\n" +
	"\aRooster\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x04R\x06userId\x12\x16\n" +
	"\x06origin\x18\x03 \x01(\tR\x06origin\x12\x12\n" +
	"\x04type\x18\x04 \x01(\x05R\x04type\"x\n" +
	"\x04Item\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x04R\x06userId\x12\x1a\n" +
	"\bquantity\x18\x03 \x01(\x05R\bquantity\x12\x17\n" +
	"\aitem_id\x18\x04 \x01(\x05R\x06itemId\x12\x12\n" +
	"\x04type\x18\x05 \x01(\x05R\x04type\"w\n" +
	"\x16GetUserProfileResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12*\n" +
	"\broosters\x18\x02 \x03(\v2\x0e.trade.RoosterR\broosters\x12!\n" +
	"\x05items\x18\x03 \x03(\v2\v.trade.ItemR\x05items\"J\n" +
	"\x13FinishTradeResponse\x12\x0e\n" +
	"\x02ok\x18\x01 \x01(\bR\x02ok\x12\x19\n" +
	"\x05error\x18\x02 \x01(\tH\x00R\x05error\x88\x01\x01B\b\n" +
	"\x06_error2\xa0\x01\n" +
	"\x05Trade\x12O\n" +
	"\x0eGetUserProfile\x12\x1c.trade.GetUserProfileRequest\x1a\x1d.trade.GetUserProfileResponse\"\x00\x12F\n" +
	"\vFinishTrade\x12\x19.trade.FinishTradeRequest\x1a\x1a.trade.FinishTradeResponse\"\x00B\tZ\a./protob\x06proto3"

var (
	file_proto_trade_proto_rawDescOnce sync.Once
	file_proto_trade_proto_rawDescData []byte
)

func file_proto_trade_proto_rawDescGZIP() []byte {
	file_proto_trade_proto_rawDescOnce.Do(func() {
		file_proto_trade_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_trade_proto_rawDesc), len(file_proto_trade_proto_rawDesc)))
	})
	return file_proto_trade_proto_rawDescData
}

var file_proto_trade_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_trade_proto_goTypes = []any{
	(*GetUserProfileRequest)(nil),  // 0: trade.GetUserProfileRequest
	(*TradeItem)(nil),              // 1: trade.TradeItem
	(*FinishTradeRequest)(nil),     // 2: trade.FinishTradeRequest
	(*Rooster)(nil),                // 3: trade.Rooster
	(*Item)(nil),                   // 4: trade.Item
	(*GetUserProfileResponse)(nil), // 5: trade.GetUserProfileResponse
	(*FinishTradeResponse)(nil),    // 6: trade.FinishTradeResponse
}
var file_proto_trade_proto_depIdxs = []int32{
	1, // 0: trade.FinishTradeRequest.items:type_name -> trade.TradeItem
	3, // 1: trade.GetUserProfileResponse.roosters:type_name -> trade.Rooster
	4, // 2: trade.GetUserProfileResponse.items:type_name -> trade.Item
	0, // 3: trade.Trade.GetUserProfile:input_type -> trade.GetUserProfileRequest
	2, // 4: trade.Trade.FinishTrade:input_type -> trade.FinishTradeRequest
	5, // 5: trade.Trade.GetUserProfile:output_type -> trade.GetUserProfileResponse
	6, // 6: trade.Trade.FinishTrade:output_type -> trade.FinishTradeResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_trade_proto_init() }
func file_proto_trade_proto_init() {
	if File_proto_trade_proto != nil {
		return
	}
	file_proto_trade_proto_msgTypes[6].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_trade_proto_rawDesc), len(file_proto_trade_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_trade_proto_goTypes,
		DependencyIndexes: file_proto_trade_proto_depIdxs,
		MessageInfos:      file_proto_trade_proto_msgTypes,
	}.Build()
	File_proto_trade_proto = out.File
	file_proto_trade_proto_goTypes = nil
	file_proto_trade_proto_depIdxs = nil
}
