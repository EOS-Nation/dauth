// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dfuse/billing/v1/billing.proto

package pbbilling

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Command struct {
	// Types that are valid to be assigned to Action:
	//	*Command_BlackListUserAction
	//	*Command_UnBlackListUserAction
	//	*Command_EventAction
	Action               isCommand_Action `protobuf_oneof:"action"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_6842c46c6c3c8e29, []int{0}
}

func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (m *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(m, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

type isCommand_Action interface {
	isCommand_Action()
}

type Command_BlackListUserAction struct {
	BlackListUserAction *BlackListUserAction `protobuf:"bytes,1,opt,name=black_list_user_action,json=blackListUserAction,proto3,oneof"`
}

type Command_UnBlackListUserAction struct {
	UnBlackListUserAction *UnBlackListUserAction `protobuf:"bytes,2,opt,name=un_black_list_user_action,json=unBlackListUserAction,proto3,oneof"`
}

type Command_EventAction struct {
	EventAction *EventAction `protobuf:"bytes,3,opt,name=event_action,json=eventAction,proto3,oneof"`
}

func (*Command_BlackListUserAction) isCommand_Action() {}

func (*Command_UnBlackListUserAction) isCommand_Action() {}

func (*Command_EventAction) isCommand_Action() {}

func (m *Command) GetAction() isCommand_Action {
	if m != nil {
		return m.Action
	}
	return nil
}

func (m *Command) GetBlackListUserAction() *BlackListUserAction {
	if x, ok := m.GetAction().(*Command_BlackListUserAction); ok {
		return x.BlackListUserAction
	}
	return nil
}

func (m *Command) GetUnBlackListUserAction() *UnBlackListUserAction {
	if x, ok := m.GetAction().(*Command_UnBlackListUserAction); ok {
		return x.UnBlackListUserAction
	}
	return nil
}

func (m *Command) GetEventAction() *EventAction {
	if x, ok := m.GetAction().(*Command_EventAction); ok {
		return x.EventAction
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Command) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Command_BlackListUserAction)(nil),
		(*Command_UnBlackListUserAction)(nil),
		(*Command_EventAction)(nil),
	}
}

type BlackListUserAction struct {
	UserId               string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Reason               string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	TtlMillis            int64    `protobuf:"varint,3,opt,name=ttl_millis,json=ttlMillis,proto3" json:"ttl_millis,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlackListUserAction) Reset()         { *m = BlackListUserAction{} }
func (m *BlackListUserAction) String() string { return proto.CompactTextString(m) }
func (*BlackListUserAction) ProtoMessage()    {}
func (*BlackListUserAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_6842c46c6c3c8e29, []int{1}
}

func (m *BlackListUserAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlackListUserAction.Unmarshal(m, b)
}
func (m *BlackListUserAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlackListUserAction.Marshal(b, m, deterministic)
}
func (m *BlackListUserAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlackListUserAction.Merge(m, src)
}
func (m *BlackListUserAction) XXX_Size() int {
	return xxx_messageInfo_BlackListUserAction.Size(m)
}
func (m *BlackListUserAction) XXX_DiscardUnknown() {
	xxx_messageInfo_BlackListUserAction.DiscardUnknown(m)
}

var xxx_messageInfo_BlackListUserAction proto.InternalMessageInfo

func (m *BlackListUserAction) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *BlackListUserAction) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *BlackListUserAction) GetTtlMillis() int64 {
	if m != nil {
		return m.TtlMillis
	}
	return 0
}

type UnBlackListUserAction struct {
	UserId               string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnBlackListUserAction) Reset()         { *m = UnBlackListUserAction{} }
func (m *UnBlackListUserAction) String() string { return proto.CompactTextString(m) }
func (*UnBlackListUserAction) ProtoMessage()    {}
func (*UnBlackListUserAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_6842c46c6c3c8e29, []int{2}
}

func (m *UnBlackListUserAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnBlackListUserAction.Unmarshal(m, b)
}
func (m *UnBlackListUserAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnBlackListUserAction.Marshal(b, m, deterministic)
}
func (m *UnBlackListUserAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnBlackListUserAction.Merge(m, src)
}
func (m *UnBlackListUserAction) XXX_Size() int {
	return xxx_messageInfo_UnBlackListUserAction.Size(m)
}
func (m *UnBlackListUserAction) XXX_DiscardUnknown() {
	xxx_messageInfo_UnBlackListUserAction.DiscardUnknown(m)
}

var xxx_messageInfo_UnBlackListUserAction proto.InternalMessageInfo

func (m *UnBlackListUserAction) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type EventAction struct {
	Event                *Event   `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventAction) Reset()         { *m = EventAction{} }
func (m *EventAction) String() string { return proto.CompactTextString(m) }
func (*EventAction) ProtoMessage()    {}
func (*EventAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_6842c46c6c3c8e29, []int{3}
}

func (m *EventAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventAction.Unmarshal(m, b)
}
func (m *EventAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventAction.Marshal(b, m, deterministic)
}
func (m *EventAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventAction.Merge(m, src)
}
func (m *EventAction) XXX_Size() int {
	return xxx_messageInfo_EventAction.Size(m)
}
func (m *EventAction) XXX_DiscardUnknown() {
	xxx_messageInfo_EventAction.DiscardUnknown(m)
}

var xxx_messageInfo_EventAction proto.InternalMessageInfo

func (m *EventAction) GetEvent() *Event {
	if m != nil {
		return m.Event
	}
	return nil
}

type Event struct {
	UserId               string               `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ApiKeyId             string               `protobuf:"bytes,2,opt,name=api_key_id,json=apiKeyId,proto3" json:"api_key_id,omitempty"`
	Usage                string               `protobuf:"bytes,3,opt,name=usage,proto3" json:"usage,omitempty"`
	Source               string               `protobuf:"bytes,4,opt,name=source,proto3" json:"source,omitempty"`
	IpAddress            string               `protobuf:"bytes,5,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	Kind                 string               `protobuf:"bytes,6,opt,name=kind,proto3" json:"kind,omitempty"`
	Network              string               `protobuf:"bytes,7,opt,name=network,proto3" json:"network,omitempty"`
	Method               string               `protobuf:"bytes,8,opt,name=method,proto3" json:"method,omitempty"`
	RequestsCount        int64                `protobuf:"varint,20,opt,name=requests_count,json=requestsCount,proto3" json:"requests_count,omitempty"`
	ResponsesCount       int64                `protobuf:"varint,21,opt,name=responses_count,json=responsesCount,proto3" json:"responses_count,omitempty"`
	RateLimitHitCount    int64                `protobuf:"varint,22,opt,name=rate_limit_hit_count,json=rateLimitHitCount,proto3" json:"rate_limit_hit_count,omitempty"`
	IngressBytes         int64                `protobuf:"varint,23,opt,name=ingress_bytes,json=ingressBytes,proto3" json:"ingress_bytes,omitempty"`
	EgressBytes          int64                `protobuf:"varint,24,opt,name=egress_bytes,json=egressBytes,proto3" json:"egress_bytes,omitempty"`
	IdleTime             int64                `protobuf:"varint,25,opt,name=idle_time,json=idleTime,proto3" json:"idle_time,omitempty"`
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,26,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_6842c46c6c3c8e29, []int{4}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Event) GetApiKeyId() string {
	if m != nil {
		return m.ApiKeyId
	}
	return ""
}

func (m *Event) GetUsage() string {
	if m != nil {
		return m.Usage
	}
	return ""
}

func (m *Event) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Event) GetIpAddress() string {
	if m != nil {
		return m.IpAddress
	}
	return ""
}

func (m *Event) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *Event) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

func (m *Event) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Event) GetRequestsCount() int64 {
	if m != nil {
		return m.RequestsCount
	}
	return 0
}

func (m *Event) GetResponsesCount() int64 {
	if m != nil {
		return m.ResponsesCount
	}
	return 0
}

func (m *Event) GetRateLimitHitCount() int64 {
	if m != nil {
		return m.RateLimitHitCount
	}
	return 0
}

func (m *Event) GetIngressBytes() int64 {
	if m != nil {
		return m.IngressBytes
	}
	return 0
}

func (m *Event) GetEgressBytes() int64 {
	if m != nil {
		return m.EgressBytes
	}
	return 0
}

func (m *Event) GetIdleTime() int64 {
	if m != nil {
		return m.IdleTime
	}
	return 0
}

func (m *Event) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func init() {
	proto.RegisterType((*Command)(nil), "dfuse.billing.v1.Command")
	proto.RegisterType((*BlackListUserAction)(nil), "dfuse.billing.v1.BlackListUserAction")
	proto.RegisterType((*UnBlackListUserAction)(nil), "dfuse.billing.v1.UnBlackListUserAction")
	proto.RegisterType((*EventAction)(nil), "dfuse.billing.v1.EventAction")
	proto.RegisterType((*Event)(nil), "dfuse.billing.v1.Event")
}

func init() { proto.RegisterFile("dfuse/billing/v1/billing.proto", fileDescriptor_6842c46c6c3c8e29) }

var fileDescriptor_6842c46c6c3c8e29 = []byte{
	// 575 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0x5b, 0x4f, 0xd4, 0x40,
	0x14, 0x96, 0xcb, 0x5e, 0x7a, 0x16, 0x50, 0x87, 0x5b, 0x41, 0xf1, 0xb2, 0x86, 0xe0, 0x0b, 0xad,
	0xe8, 0x8b, 0x89, 0xbe, 0xb0, 0xc4, 0x04, 0x22, 0xbe, 0x34, 0xf0, 0x62, 0x4c, 0x26, 0xd3, 0xed,
	0xa1, 0x4c, 0xb6, 0xed, 0xd4, 0xce, 0x74, 0x0d, 0x7f, 0xc5, 0x7f, 0xe9, 0x3f, 0x30, 0x73, 0xa6,
	0x05, 0x02, 0xbb, 0x6f, 0xf3, 0x5d, 0xe6, 0x7c, 0xfd, 0xda, 0x3d, 0x0b, 0xaf, 0x92, 0xab, 0x5a,
	0x63, 0x18, 0xcb, 0x2c, 0x93, 0x45, 0x1a, 0x4e, 0x8f, 0xda, 0x63, 0x50, 0x56, 0xca, 0x28, 0xf6,
	0x8c, 0xf4, 0xa0, 0x25, 0xa7, 0x47, 0xbb, 0xaf, 0x53, 0xa5, 0xd2, 0x0c, 0x43, 0xd2, 0xe3, 0xfa,
	0x2a, 0x34, 0x32, 0x47, 0x6d, 0x44, 0x5e, 0xba, 0x2b, 0xc3, 0xbf, 0x8b, 0xd0, 0x3b, 0x51, 0x79,
	0x2e, 0x8a, 0x84, 0xfd, 0x82, 0xad, 0x38, 0x13, 0xe3, 0x09, 0xcf, 0xa4, 0x36, 0xbc, 0xd6, 0x58,
	0x71, 0x31, 0x36, 0x52, 0x15, 0xfe, 0xc2, 0x9b, 0x85, 0xf7, 0x83, 0x8f, 0xfb, 0xc1, 0xc3, 0xf9,
	0xc1, 0xc8, 0xfa, 0xcf, 0xa5, 0x36, 0x97, 0x1a, 0xab, 0x63, 0x32, 0x9f, 0x3e, 0x89, 0xd6, 0xe3,
	0xc7, 0x34, 0x1b, 0xc3, 0x4e, 0x5d, 0xf0, 0x39, 0x01, 0x8b, 0x14, 0x70, 0xf0, 0x38, 0xe0, 0xb2,
	0x98, 0x1d, 0xb1, 0x59, 0xcf, 0x12, 0xd8, 0x08, 0x56, 0x70, 0x8a, 0x85, 0x69, 0xe7, 0x2e, 0xd1,
	0xdc, 0xbd, 0xc7, 0x73, 0xbf, 0x59, 0xd7, 0xed, 0xb4, 0x01, 0xde, 0xc1, 0x51, 0x1f, 0xba, 0xee,
	0xf6, 0x10, 0x61, 0x7d, 0x56, 0xc8, 0x36, 0xf4, 0xe8, 0xd9, 0x65, 0x42, 0x2f, 0xc6, 0x8b, 0xba,
	0x16, 0x9e, 0x25, 0x6c, 0x0b, 0xba, 0x15, 0x0a, 0xdd, 0xf4, 0xf1, 0xa2, 0x06, 0xb1, 0x3d, 0x00,
	0x63, 0x32, 0x9e, 0xdb, 0x78, 0x4d, 0xcf, 0xb4, 0x14, 0x79, 0xc6, 0x64, 0x3f, 0x88, 0x18, 0x7e,
	0x80, 0xcd, 0x99, 0x35, 0xe7, 0x06, 0x0d, 0xbf, 0xc2, 0xe0, 0x5e, 0x01, 0x76, 0x08, 0x1d, 0x2a,
	0xd0, 0x7c, 0xa7, 0xed, 0x39, 0x75, 0x23, 0xe7, 0x1a, 0xfe, 0x5b, 0x82, 0x0e, 0x11, 0xf3, 0x9b,
	0xbc, 0x04, 0x10, 0xa5, 0xe4, 0x13, 0xbc, 0xb1, 0x9a, 0x6b, 0xd3, 0x17, 0xa5, 0xfc, 0x8e, 0x37,
	0x67, 0x09, 0xdb, 0x80, 0x4e, 0xad, 0x45, 0x8a, 0x54, 0xc5, 0x8b, 0x1c, 0xb0, 0xed, 0xb5, 0xaa,
	0xab, 0x31, 0xfa, 0xcb, 0x6e, 0x96, 0x43, 0xb6, 0xbd, 0x2c, 0xb9, 0x48, 0x92, 0x0a, 0xb5, 0xf6,
	0x3b, 0xa4, 0x79, 0xb2, 0x3c, 0x76, 0x04, 0x63, 0xb0, 0x3c, 0x91, 0x45, 0xe2, 0x77, 0x49, 0xa0,
	0x33, 0xf3, 0xa1, 0x57, 0xa0, 0xf9, 0xa3, 0xaa, 0x89, 0xdf, 0x23, 0xba, 0x85, 0x36, 0x24, 0x47,
	0x73, 0xad, 0x12, 0xbf, 0xef, 0x42, 0x1c, 0x62, 0xfb, 0xb0, 0x56, 0xe1, 0xef, 0x1a, 0xb5, 0xd1,
	0x7c, 0xac, 0xea, 0xc2, 0xf8, 0x1b, 0xf4, 0x9a, 0x57, 0x5b, 0xf6, 0xc4, 0x92, 0xec, 0x00, 0x9e,
	0x56, 0xa8, 0x4b, 0x55, 0x68, 0x6c, 0x7d, 0x9b, 0xe4, 0x5b, 0xbb, 0xa5, 0x9d, 0x31, 0x84, 0x8d,
	0x4a, 0x18, 0xe4, 0x99, 0xcc, 0xa5, 0xe1, 0xd7, 0xd2, 0x34, 0xee, 0x2d, 0x72, 0x3f, 0xb7, 0xda,
	0xb9, 0x95, 0x4e, 0xa5, 0x71, 0x17, 0xde, 0xc1, 0xaa, 0x2c, 0x52, 0xdb, 0x88, 0xc7, 0x37, 0x06,
	0xb5, 0xbf, 0x4d, 0xce, 0x95, 0x86, 0x1c, 0x59, 0x8e, 0xbd, 0x85, 0x15, 0xbc, 0xef, 0xf1, 0xc9,
	0x33, 0xc0, 0x7b, 0x96, 0x17, 0xe0, 0xc9, 0x24, 0x43, 0x6e, 0x17, 0xd5, 0xdf, 0x21, 0xbd, 0x6f,
	0x89, 0x0b, 0x99, 0x23, 0xfb, 0x0c, 0xde, 0xed, 0x02, 0xfb, 0xbb, 0xf4, 0xb1, 0x77, 0x03, 0xb7,
	0xe2, 0x41, 0xbb, 0xe2, 0xc1, 0x45, 0xeb, 0x88, 0xee, 0xcc, 0xa3, 0x93, 0x9f, 0xc7, 0xa9, 0x34,
	0xd7, 0x75, 0x1c, 0x8c, 0x55, 0x1e, 0xd2, 0xef, 0xe3, 0x50, 0xaa, 0xe6, 0xa0, 0x85, 0xd0, 0x87,
	0x65, 0x25, 0xa7, 0x61, 0x19, 0x87, 0x0f, 0xff, 0x63, 0xbe, 0x94, 0x71, 0x03, 0xe2, 0x2e, 0x65,
	0x7c, 0xfa, 0x1f, 0x00, 0x00, 0xff, 0xff, 0xce, 0x49, 0x83, 0xcd, 0x88, 0x04, 0x00, 0x00,
}
