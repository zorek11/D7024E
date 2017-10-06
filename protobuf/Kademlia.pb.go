// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Kademlia.proto

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	Kademlia.proto

It has these top-level messages:
	KademliaMessage
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type KademliaMessage struct {
	Label            *string                        `protobuf:"bytes,1,req,name=label" json:"label,omitempty"`
	Senderid         *string                        `protobuf:"bytes,2,req,name=senderid" json:"senderid,omitempty"`
	SenderAddr       *string                        `protobuf:"bytes,3,req,name=senderAddr" json:"senderAddr,omitempty"`
	Lookupcontact    *KademliaMessage_LookupContact `protobuf:"group,4,opt,name=LookupContact" json:"lookupcontact,omitempty"`
	Data             []byte                         `protobuf:"bytes,8,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte                         `json:"-"`
}

func (m *KademliaMessage) Reset()                    { *m = KademliaMessage{} }
func (m *KademliaMessage) String() string            { return proto.CompactTextString(m) }
func (*KademliaMessage) ProtoMessage()               {}
func (*KademliaMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *KademliaMessage) GetLabel() string {
	if m != nil && m.Label != nil {
		return *m.Label
	}
	return ""
}

func (m *KademliaMessage) GetSenderid() string {
	if m != nil && m.Senderid != nil {
		return *m.Senderid
	}
	return ""
}

func (m *KademliaMessage) GetSenderAddr() string {
	if m != nil && m.SenderAddr != nil {
		return *m.SenderAddr
	}
	return ""
}

func (m *KademliaMessage) GetLookupcontact() *KademliaMessage_LookupContact {
	if m != nil {
		return m.Lookupcontact
	}
	return nil
}

func (m *KademliaMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type KademliaMessage_LookupContact struct {
	Id               *string `protobuf:"bytes,5,req,name=id" json:"id,omitempty"`
	Address          *string `protobuf:"bytes,6,opt,name=address" json:"address,omitempty"`
	Distance         *string `protobuf:"bytes,7,opt,name=distance" json:"distance,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *KademliaMessage_LookupContact) Reset()         { *m = KademliaMessage_LookupContact{} }
func (m *KademliaMessage_LookupContact) String() string { return proto.CompactTextString(m) }
func (*KademliaMessage_LookupContact) ProtoMessage()    {}
func (*KademliaMessage_LookupContact) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 0}
}

func (m *KademliaMessage_LookupContact) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *KademliaMessage_LookupContact) GetAddress() string {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return ""
}

func (m *KademliaMessage_LookupContact) GetDistance() string {
	if m != nil && m.Distance != nil {
		return *m.Distance
	}
	return ""
}

func init() {
	proto.RegisterType((*KademliaMessage)(nil), "protobuf.KademliaMessage")
	proto.RegisterType((*KademliaMessage_LookupContact)(nil), "protobuf.KademliaMessage.LookupContact")
}

func init() { proto.RegisterFile("Kademlia.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x8e, 0xb1, 0x0e, 0x82, 0x30,
	0x14, 0x45, 0x03, 0x82, 0xe0, 0x0b, 0x88, 0xe9, 0xd4, 0x38, 0x11, 0x17, 0x99, 0x3a, 0xf8, 0x01,
	0x26, 0xc6, 0x51, 0xfd, 0x88, 0x07, 0x7d, 0x1a, 0x62, 0xa5, 0x84, 0x96, 0x3f, 0xf5, 0x83, 0xac,
	0x25, 0x0c, 0x3a, 0x35, 0x39, 0xf7, 0xf6, 0x9e, 0x07, 0xeb, 0x0b, 0x4a, 0x7a, 0xa9, 0x16, 0x45,
	0x3f, 0x68, 0xab, 0x59, 0xea, 0x9f, 0x7a, 0xbc, 0xef, 0xde, 0x01, 0x14, 0x73, 0x78, 0x23, 0x63,
	0xf0, 0x41, 0x2c, 0x87, 0x58, 0x61, 0x4d, 0x8a, 0x07, 0x65, 0x58, 0xad, 0xd8, 0x06, 0x52, 0x43,
	0x9d, 0xa4, 0xa1, 0x95, 0x3c, 0xf4, 0x84, 0x01, 0x4c, 0xe4, 0x24, 0xe5, 0xc0, 0x17, 0x9e, 0x1d,
	0x21, 0x57, 0x5a, 0x3f, 0xc7, 0xbe, 0xd1, 0x9d, 0xc5, 0xc6, 0xf2, 0xa8, 0x0c, 0x2a, 0x38, 0xec,
	0xc5, 0xac, 0x12, 0x7f, 0x1a, 0x71, 0xf5, 0xf5, 0xf3, 0x54, 0x67, 0x19, 0x44, 0x12, 0x2d, 0xf2,
	0xd4, 0x7d, 0xcb, 0xb6, 0x6e, 0xed, 0x37, 0x06, 0x08, 0x9d, 0x3e, 0xf6, 0xaa, 0x02, 0x12, 0x74,
	0x62, 0x37, 0xc3, 0x97, 0xae, 0xed, 0x2f, 0x94, 0xad, 0xb1, 0xd8, 0x35, 0xc4, 0x93, 0x2f, 0xf9,
	0x04, 0x00, 0x00, 0xff, 0xff, 0xae, 0x38, 0xe6, 0x20, 0xf1, 0x00, 0x00, 0x00,
}
