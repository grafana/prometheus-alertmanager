// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: silence.proto

package silencepb

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Type specifies how the given name and pattern are matched
// against a label set.
type Matcher_Type int32

const (
	Matcher_EQUAL      Matcher_Type = 0
	Matcher_REGEXP     Matcher_Type = 1
	Matcher_NOT_EQUAL  Matcher_Type = 2
	Matcher_NOT_REGEXP Matcher_Type = 3
)

var Matcher_Type_name = map[int32]string{
	0: "EQUAL",
	1: "REGEXP",
	2: "NOT_EQUAL",
	3: "NOT_REGEXP",
}

var Matcher_Type_value = map[string]int32{
	"EQUAL":      0,
	"REGEXP":     1,
	"NOT_EQUAL":  2,
	"NOT_REGEXP": 3,
}

func (x Matcher_Type) String() string {
	return proto.EnumName(Matcher_Type_name, int32(x))
}

func (Matcher_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7fc56058cf68dbd8, []int{0, 0}
}

// Matcher specifies a rule, which can match or set of labels or not.
type Matcher struct {
	Type Matcher_Type `protobuf:"varint,1,opt,name=type,proto3,enum=silencepb.Matcher_Type" json:"type,omitempty"`
	// The label name in a label set to against which the matcher
	// checks the pattern.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The pattern being checked according to the matcher's type.
	Pattern              string   `protobuf:"bytes,3,opt,name=pattern,proto3" json:"pattern,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Matcher) Reset()         { *m = Matcher{} }
func (m *Matcher) String() string { return proto.CompactTextString(m) }
func (*Matcher) ProtoMessage()    {}
func (*Matcher) Descriptor() ([]byte, []int) {
	return fileDescriptor_7fc56058cf68dbd8, []int{0}
}
func (m *Matcher) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Matcher) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Matcher.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Matcher) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Matcher.Merge(m, src)
}
func (m *Matcher) XXX_Size() int {
	return m.Size()
}
func (m *Matcher) XXX_DiscardUnknown() {
	xxx_messageInfo_Matcher.DiscardUnknown(m)
}

var xxx_messageInfo_Matcher proto.InternalMessageInfo

// DEPRECATED: A comment can be attached to a silence.
type Comment struct {
	Author               string    `protobuf:"bytes,1,opt,name=author,proto3" json:"author,omitempty"`
	Comment              string    `protobuf:"bytes,2,opt,name=comment,proto3" json:"comment,omitempty"`
	Timestamp            time.Time `protobuf:"bytes,3,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Comment) Reset()         { *m = Comment{} }
func (m *Comment) String() string { return proto.CompactTextString(m) }
func (*Comment) ProtoMessage()    {}
func (*Comment) Descriptor() ([]byte, []int) {
	return fileDescriptor_7fc56058cf68dbd8, []int{1}
}
func (m *Comment) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Comment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Comment.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Comment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Comment.Merge(m, src)
}
func (m *Comment) XXX_Size() int {
	return m.Size()
}
func (m *Comment) XXX_DiscardUnknown() {
	xxx_messageInfo_Comment.DiscardUnknown(m)
}

var xxx_messageInfo_Comment proto.InternalMessageInfo

// Silence specifies an object that ignores alerts based
// on a set of matchers during a given time frame.
type Silence struct {
	// A globally unique identifier.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// A set of matchers all of which have to be true for a silence
	// to affect a given label set.
	Matchers []*Matcher `protobuf:"bytes,2,rep,name=matchers,proto3" json:"matchers,omitempty"`
	// The time range during which the silence is active.
	StartsAt time.Time `protobuf:"bytes,3,opt,name=starts_at,json=startsAt,proto3,stdtime" json:"starts_at"`
	EndsAt   time.Time `protobuf:"bytes,4,opt,name=ends_at,json=endsAt,proto3,stdtime" json:"ends_at"`
	// The last notification made to the silence.
	UpdatedAt time.Time `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3,stdtime" json:"updated_at"`
	// DEPRECATED: A set of comments made on the silence.
	Comments []*Comment `protobuf:"bytes,7,rep,name=comments,proto3" json:"comments,omitempty"`
	// Comment for the silence.
	CreatedBy            string   `protobuf:"bytes,8,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	Comment              string   `protobuf:"bytes,9,opt,name=comment,proto3" json:"comment,omitempty"`
	Name                 string   `json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Silence) Reset()         { *m = Silence{} }
func (m *Silence) String() string { return proto.CompactTextString(m) }
func (*Silence) ProtoMessage()    {}
func (*Silence) Descriptor() ([]byte, []int) {
	return fileDescriptor_7fc56058cf68dbd8, []int{2}
}
func (m *Silence) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Silence) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Silence.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Silence) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Silence.Merge(m, src)
}
func (m *Silence) XXX_Size() int {
	return m.Size()
}
func (m *Silence) XXX_DiscardUnknown() {
	xxx_messageInfo_Silence.DiscardUnknown(m)
}

var xxx_messageInfo_Silence proto.InternalMessageInfo

// MeshSilence wraps a regular silence with an expiration timestamp
// after which the silence may be garbage collected.
type MeshSilence struct {
	Silence              *Silence  `protobuf:"bytes,1,opt,name=silence,proto3" json:"silence,omitempty"`
	ExpiresAt            time.Time `protobuf:"bytes,2,opt,name=expires_at,json=expiresAt,proto3,stdtime" json:"expires_at"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *MeshSilence) Reset()         { *m = MeshSilence{} }
func (m *MeshSilence) String() string { return proto.CompactTextString(m) }
func (*MeshSilence) ProtoMessage()    {}
func (*MeshSilence) Descriptor() ([]byte, []int) {
	return fileDescriptor_7fc56058cf68dbd8, []int{3}
}
func (m *MeshSilence) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MeshSilence) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MeshSilence.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MeshSilence) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MeshSilence.Merge(m, src)
}
func (m *MeshSilence) XXX_Size() int {
	return m.Size()
}
func (m *MeshSilence) XXX_DiscardUnknown() {
	xxx_messageInfo_MeshSilence.DiscardUnknown(m)
}

var xxx_messageInfo_MeshSilence proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("silencepb.Matcher_Type", Matcher_Type_name, Matcher_Type_value)
	proto.RegisterType((*Matcher)(nil), "silencepb.Matcher")
	proto.RegisterType((*Comment)(nil), "silencepb.Comment")
	proto.RegisterType((*Silence)(nil), "silencepb.Silence")
	proto.RegisterType((*MeshSilence)(nil), "silencepb.MeshSilence")
}

func init() { proto.RegisterFile("silence.proto", fileDescriptor_7fc56058cf68dbd8) }

var fileDescriptor_7fc56058cf68dbd8 = []byte{
	// 461 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcd, 0x8e, 0xd2, 0x50,
	0x14, 0xc7, 0xb9, 0x85, 0xa1, 0xdc, 0x43, 0x86, 0x90, 0x13, 0xa3, 0x0d, 0x89, 0x40, 0xba, 0x22,
	0xd1, 0x94, 0x04, 0xb7, 0xba, 0x28, 0x13, 0xe2, 0xc6, 0xf1, 0xa3, 0x62, 0xe2, 0x6e, 0x52, 0xe8,
	0x11, 0x9a, 0x4c, 0x3f, 0xd2, 0x1e, 0x12, 0x59, 0xe9, 0x23, 0xf8, 0x0c, 0x3e, 0x0d, 0x4b, 0x9f,
	0xc0, 0x0f, 0xde, 0xc2, 0x9d, 0xe9, 0xed, 0x2d, 0xce, 0x84, 0x55, 0x77, 0xf7, 0xdc, 0xf3, 0xff,
	0x9f, 0x8f, 0xdf, 0x81, 0xcb, 0x3c, 0xbc, 0xa5, 0x78, 0x4d, 0x4e, 0x9a, 0x25, 0x9c, 0xa0, 0xd4,
	0x61, 0xba, 0x1a, 0x8c, 0x36, 0x49, 0xb2, 0xb9, 0xa5, 0xa9, 0x4a, 0xac, 0x76, 0x9f, 0xa6, 0x1c,
	0x46, 0x94, 0xb3, 0x1f, 0xa5, 0xa5, 0x76, 0xf0, 0x60, 0x93, 0x6c, 0x12, 0xf5, 0x9c, 0x16, 0xaf,
	0xf2, 0xd7, 0xfe, 0x2e, 0xc0, 0xbc, 0xf6, 0x79, 0xbd, 0xa5, 0x0c, 0x9f, 0x40, 0x8b, 0xf7, 0x29,
	0x59, 0x62, 0x2c, 0x26, 0xbd, 0xd9, 0x23, 0xe7, 0x54, 0xdc, 0xd1, 0x0a, 0x67, 0xb9, 0x4f, 0xc9,
	0x53, 0x22, 0x44, 0x68, 0xc5, 0x7e, 0x44, 0x96, 0x31, 0x16, 0x13, 0xe9, 0xa9, 0x37, 0x5a, 0x60,
	0xa6, 0x3e, 0x33, 0x65, 0xb1, 0xd5, 0x54, 0xdf, 0x55, 0x68, 0x3f, 0x87, 0x56, 0xe1, 0x45, 0x09,
	0x17, 0x8b, 0x77, 0x1f, 0xdc, 0x57, 0xfd, 0x06, 0x02, 0xb4, 0xbd, 0xc5, 0xcb, 0xc5, 0xc7, 0xb7,
	0x7d, 0x81, 0x97, 0x20, 0x5f, 0xbf, 0x59, 0xde, 0x94, 0x29, 0x03, 0x7b, 0x00, 0x45, 0xa8, 0xd3,
	0x4d, 0xfb, 0x0b, 0x98, 0x57, 0x49, 0x14, 0x51, 0xcc, 0xf8, 0x10, 0xda, 0xfe, 0x8e, 0xb7, 0x49,
	0xa6, 0xa6, 0x94, 0x9e, 0x8e, 0x8a, 0xd6, 0xeb, 0x52, 0xa2, 0x27, 0xaa, 0x42, 0x9c, 0x83, 0x3c,
	0xa1, 0x50, 0x63, 0x75, 0x67, 0x03, 0xa7, 0x84, 0xe5, 0x54, 0xb0, 0x9c, 0x65, 0xa5, 0x98, 0x77,
	0x0e, 0x3f, 0x47, 0x8d, 0x6f, 0xbf, 0x46, 0xc2, 0xfb, 0x6f, 0xb3, 0xff, 0x1a, 0x60, 0xbe, 0x2f,
	0x69, 0x60, 0x0f, 0x8c, 0x30, 0xd0, 0xdd, 0x8d, 0x30, 0x40, 0x07, 0x3a, 0x51, 0x89, 0x27, 0xb7,
	0x8c, 0x71, 0x73, 0xd2, 0x9d, 0xe1, 0x39, 0x39, 0xef, 0xa4, 0x41, 0x17, 0x64, 0xce, 0x7e, 0xc6,
	0xf9, 0x8d, 0xcf, 0xb5, 0xe6, 0xe9, 0x94, 0x36, 0x97, 0xf1, 0x05, 0x98, 0x14, 0x07, 0xaa, 0x40,
	0xab, 0x46, 0x81, 0x76, 0x61, 0x72, 0x19, 0xaf, 0x00, 0x76, 0x69, 0xe0, 0x33, 0x05, 0x45, 0x85,
	0x8b, 0x3a, 0x48, 0xb4, 0xcf, 0xe5, 0x62, 0x6d, 0x4d, 0x38, 0xb7, 0xcc, 0xb3, 0xb5, 0xf5, 0xb9,
	0xbc, 0x93, 0x06, 0x1f, 0x03, 0xac, 0x33, 0x52, 0x4d, 0x57, 0x7b, 0xab, 0xa3, 0xf0, 0x49, 0xfd,
	0x33, 0xdf, 0xdf, 0xbd, 0x9f, 0xbc, 0x77, 0x3f, 0xfb, 0xab, 0x80, 0xee, 0x35, 0xe5, 0xdb, 0x8a,
	0xff, 0x53, 0x30, 0x75, 0x1f, 0x75, 0x84, 0xfb, 0x7d, 0xb5, 0xc8, 0xab, 0x24, 0xc5, 0xae, 0xf4,
	0x39, 0x0d, 0x33, 0x52, 0xb4, 0x8c, 0x3a, 0xbb, 0x6a, 0x9f, 0xcb, 0xf3, 0xfe, 0xe1, 0xcf, 0xb0,
	0x71, 0x38, 0x0e, 0xc5, 0x8f, 0xe3, 0x50, 0xfc, 0x3e, 0x0e, 0xc5, 0xaa, 0xad, 0xac, 0xcf, 0xfe,
	0x05, 0x00, 0x00, 0xff, 0xff, 0xad, 0x09, 0x3a, 0xe9, 0x90, 0x03, 0x00, 0x00,
}

func (m *Matcher) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Matcher) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Matcher) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Pattern) > 0 {
		i -= len(m.Pattern)
		copy(dAtA[i:], m.Pattern)
		i = encodeVarintSilence(dAtA, i, uint64(len(m.Pattern)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintSilence(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if m.Type != 0 {
		i = encodeVarintSilence(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Comment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Comment) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Comment) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintSilence(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x1a
	if len(m.Comment) > 0 {
		i -= len(m.Comment)
		copy(dAtA[i:], m.Comment)
		i = encodeVarintSilence(dAtA, i, uint64(len(m.Comment)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Author) > 0 {
		i -= len(m.Author)
		copy(dAtA[i:], m.Author)
		i = encodeVarintSilence(dAtA, i, uint64(len(m.Author)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Silence) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Silence) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Silence) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Comment) > 0 {
		i -= len(m.Comment)
		copy(dAtA[i:], m.Comment)
		i = encodeVarintSilence(dAtA, i, uint64(len(m.Comment)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.CreatedBy) > 0 {
		i -= len(m.CreatedBy)
		copy(dAtA[i:], m.CreatedBy)
		i = encodeVarintSilence(dAtA, i, uint64(len(m.CreatedBy)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.Comments) > 0 {
		for iNdEx := len(m.Comments) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Comments[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSilence(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	n2, err2 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintSilence(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x2a
	n3, err3 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.EndsAt, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.EndsAt):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintSilence(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x22
	n4, err4 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.StartsAt, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.StartsAt):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintSilence(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x1a
	if len(m.Matchers) > 0 {
		for iNdEx := len(m.Matchers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Matchers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSilence(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintSilence(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MeshSilence) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MeshSilence) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MeshSilence) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	n5, err5 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.ExpiresAt, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.ExpiresAt):])
	if err5 != nil {
		return 0, err5
	}
	i -= n5
	i = encodeVarintSilence(dAtA, i, uint64(n5))
	i--
	dAtA[i] = 0x12
	if m.Silence != nil {
		{
			size, err := m.Silence.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSilence(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSilence(dAtA []byte, offset int, v uint64) int {
	offset -= sovSilence(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Matcher) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovSilence(uint64(m.Type))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovSilence(uint64(l))
	}
	l = len(m.Pattern)
	if l > 0 {
		n += 1 + l + sovSilence(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Comment) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Author)
	if l > 0 {
		n += 1 + l + sovSilence(uint64(l))
	}
	l = len(m.Comment)
	if l > 0 {
		n += 1 + l + sovSilence(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovSilence(uint64(l))
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Silence) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovSilence(uint64(l))
	}
	if len(m.Matchers) > 0 {
		for _, e := range m.Matchers {
			l = e.Size()
			n += 1 + l + sovSilence(uint64(l))
		}
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.StartsAt)
	n += 1 + l + sovSilence(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.EndsAt)
	n += 1 + l + sovSilence(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)
	n += 1 + l + sovSilence(uint64(l))
	if len(m.Comments) > 0 {
		for _, e := range m.Comments {
			l = e.Size()
			n += 1 + l + sovSilence(uint64(l))
		}
	}
	l = len(m.CreatedBy)
	if l > 0 {
		n += 1 + l + sovSilence(uint64(l))
	}
	l = len(m.Comment)
	if l > 0 {
		n += 1 + l + sovSilence(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *MeshSilence) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Silence != nil {
		l = m.Silence.Size()
		n += 1 + l + sovSilence(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.ExpiresAt)
	n += 1 + l + sovSilence(uint64(l))
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovSilence(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSilence(x uint64) (n int) {
	return sovSilence(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Matcher) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSilence
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Matcher: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Matcher: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= Matcher_Type(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pattern", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pattern = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSilence(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSilence
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Comment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSilence
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Comment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Comment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Author", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Author = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Comment", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Comment = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSilence(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSilence
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Silence) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSilence
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Silence: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Silence: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Matchers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Matchers = append(m.Matchers, &Matcher{})
			if err := m.Matchers[len(m.Matchers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartsAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.StartsAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndsAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.EndsAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Comments", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Comments = append(m.Comments, &Comment{})
			if err := m.Comments[len(m.Comments)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedBy", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CreatedBy = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Comment", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Comment = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSilence(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSilence
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MeshSilence) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSilence
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MeshSilence: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MeshSilence: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Silence", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Silence == nil {
				m.Silence = &Silence{}
			}
			if err := m.Silence.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpiresAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSilence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSilence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.ExpiresAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSilence(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSilence
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSilence(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSilence
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSilence
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSilence
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSilence
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSilence
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSilence        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSilence          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSilence = fmt.Errorf("proto: unexpected end of group")
)
