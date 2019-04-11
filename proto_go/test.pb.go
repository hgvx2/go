// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

//包名，通过protoc生成时go文件时

package test

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

//手机类型
//枚举类型第一个字段必须为0
type PhoneType int32

const (
	PhoneType_HOME PhoneType = 0
	PhoneType_WORK PhoneType = 1
)

var PhoneType_name = map[int32]string{
	0: "HOME",
	1: "WORK",
}

var PhoneType_value = map[string]int32{
	"HOME": 0,
	"WORK": 1,
}

func (x PhoneType) Enum() *PhoneType {
	p := new(PhoneType)
	*p = x
	return p
}

func (x PhoneType) String() string {
	return proto.EnumName(PhoneType_name, int32(x))
}

func (x *PhoneType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PhoneType_value, data, "PhoneType")
	if err != nil {
		return err
	}
	*x = PhoneType(value)
	return nil
}

func (PhoneType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}

//手机
type Phone struct {
	Type                 *PhoneType `protobuf:"varint,1,req,name=type,enum=test.PhoneType" json:"type,omitempty"`
	Number               *string    `protobuf:"bytes,2,req,name=number" json:"number,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Phone) Reset()         { *m = Phone{} }
func (m *Phone) String() string { return proto.CompactTextString(m) }
func (*Phone) ProtoMessage()    {}
func (*Phone) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}

func (m *Phone) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Phone.Unmarshal(m, b)
}
func (m *Phone) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Phone.Marshal(b, m, deterministic)
}
func (m *Phone) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Phone.Merge(m, src)
}
func (m *Phone) XXX_Size() int {
	return xxx_messageInfo_Phone.Size(m)
}
func (m *Phone) XXX_DiscardUnknown() {
	xxx_messageInfo_Phone.DiscardUnknown(m)
}

var xxx_messageInfo_Phone proto.InternalMessageInfo

func (m *Phone) GetType() PhoneType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return PhoneType_HOME
}

func (m *Phone) GetNumber() string {
	if m != nil && m.Number != nil {
		return *m.Number
	}
	return ""
}

//人
type Person struct {
	//后面的数字表示标识号
	Id   *int32  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Name *string `protobuf:"bytes,2,req,name=name" json:"name,omitempty"`
	//repeated表示可重复
	//可以有多个手机
	Phones               []*Phone `protobuf:"bytes,3,rep,name=phones" json:"phones,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Person) Reset()         { *m = Person{} }
func (m *Person) String() string { return proto.CompactTextString(m) }
func (*Person) ProtoMessage()    {}
func (*Person) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{1}
}

func (m *Person) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Person.Unmarshal(m, b)
}
func (m *Person) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Person.Marshal(b, m, deterministic)
}
func (m *Person) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Person.Merge(m, src)
}
func (m *Person) XXX_Size() int {
	return xxx_messageInfo_Person.Size(m)
}
func (m *Person) XXX_DiscardUnknown() {
	xxx_messageInfo_Person.DiscardUnknown(m)
}

var xxx_messageInfo_Person proto.InternalMessageInfo

func (m *Person) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Person) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Person) GetPhones() []*Phone {
	if m != nil {
		return m.Phones
	}
	return nil
}

//联系簿
type ContactBook struct {
	Persons              []*Person `protobuf:"bytes,1,rep,name=persons" json:"persons,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ContactBook) Reset()         { *m = ContactBook{} }
func (m *ContactBook) String() string { return proto.CompactTextString(m) }
func (*ContactBook) ProtoMessage()    {}
func (*ContactBook) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{2}
}

func (m *ContactBook) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContactBook.Unmarshal(m, b)
}
func (m *ContactBook) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContactBook.Marshal(b, m, deterministic)
}
func (m *ContactBook) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContactBook.Merge(m, src)
}
func (m *ContactBook) XXX_Size() int {
	return xxx_messageInfo_ContactBook.Size(m)
}
func (m *ContactBook) XXX_DiscardUnknown() {
	xxx_messageInfo_ContactBook.DiscardUnknown(m)
}

var xxx_messageInfo_ContactBook proto.InternalMessageInfo

func (m *ContactBook) GetPersons() []*Person {
	if m != nil {
		return m.Persons
	}
	return nil
}

func init() {
	proto.RegisterEnum("test.PhoneType", PhoneType_name, PhoneType_value)
	proto.RegisterType((*Phone)(nil), "test.Phone")
	proto.RegisterType((*Person)(nil), "test.Person")
	proto.RegisterType((*ContactBook)(nil), "test.ContactBook")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0xcc, 0xb8, 0x58, 0x03, 0x32,
	0xf2, 0xf3, 0x52, 0x85, 0x64, 0xb9, 0x58, 0x4a, 0x2a, 0x0b, 0x52, 0x25, 0x18, 0x15, 0x98, 0x34,
	0xf8, 0x8c, 0xf8, 0xf5, 0xc0, 0x2a, 0xc1, 0x52, 0x21, 0x95, 0x05, 0xa9, 0x42, 0x7c, 0x5c, 0x6c,
	0x79, 0xa5, 0xb9, 0x49, 0xa9, 0x45, 0x12, 0x4c, 0x0a, 0x4c, 0x1a, 0x9c, 0x4a, 0xf6, 0x5c, 0x6c,
	0x01, 0xa9, 0x45, 0xc5, 0xf9, 0x79, 0x42, 0x5c, 0x5c, 0x4c, 0x99, 0x29, 0x60, 0x6d, 0xac, 0x42,
	0x3c, 0x5c, 0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x10, 0x35, 0x42, 0xd2, 0x5c, 0x6c, 0x05, 0x20, 0x03,
	0x8a, 0x25, 0x98, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0xb8, 0x91, 0x0c, 0x55, 0xd2, 0xe1, 0xe2, 0x76,
	0xce, 0xcf, 0x2b, 0x49, 0x4c, 0x2e, 0x71, 0xca, 0xcf, 0xcf, 0x16, 0x92, 0xe5, 0x62, 0x2f, 0x00,
	0x9b, 0x57, 0x2c, 0xc1, 0x08, 0x56, 0xcc, 0x03, 0x55, 0x0c, 0x16, 0xd4, 0x92, 0xe7, 0xe2, 0x44,
	0xb8, 0x85, 0x83, 0x8b, 0xc5, 0xc3, 0xdf, 0xd7, 0x55, 0x80, 0x01, 0xc4, 0x0a, 0xf7, 0x0f, 0xf2,
	0x16, 0x60, 0x04, 0x04, 0x00, 0x00, 0xff, 0xff, 0x34, 0xed, 0x93, 0x98, 0xda, 0x00, 0x00, 0x00,
}
