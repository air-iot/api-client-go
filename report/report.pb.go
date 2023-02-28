// Code generated by protoc-gen-go. DO NOT EDIT.
// source: report/report.proto

/*
Package report is a generated protocol buffer package.

It is generated from these files:
	report/report.proto

It has these top-level messages:
	QueryReportRequest
	GetOrDeleteReportRequest
*/
package report

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/air-iot/api-client-go/v4/api"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type QueryReportRequest struct {
	Query   []byte `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Archive string `protobuf:"bytes,2,opt,name=archive" json:"archive,omitempty"`
}

func (m *QueryReportRequest) Reset()                    { *m = QueryReportRequest{} }
func (m *QueryReportRequest) String() string            { return proto.CompactTextString(m) }
func (*QueryReportRequest) ProtoMessage()               {}
func (*QueryReportRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *QueryReportRequest) GetQuery() []byte {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *QueryReportRequest) GetArchive() string {
	if m != nil {
		return m.Archive
	}
	return ""
}

type GetOrDeleteReportRequest struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Archive string `protobuf:"bytes,2,opt,name=archive" json:"archive,omitempty"`
}

func (m *GetOrDeleteReportRequest) Reset()                    { *m = GetOrDeleteReportRequest{} }
func (m *GetOrDeleteReportRequest) String() string            { return proto.CompactTextString(m) }
func (*GetOrDeleteReportRequest) ProtoMessage()               {}
func (*GetOrDeleteReportRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetOrDeleteReportRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GetOrDeleteReportRequest) GetArchive() string {
	if m != nil {
		return m.Archive
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryReportRequest)(nil), "report.QueryReportRequest")
	proto.RegisterType((*GetOrDeleteReportRequest)(nil), "report.GetOrDeleteReportRequest")
}

func init() { proto.RegisterFile("report/report.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x52, 0x41, 0x4b, 0x33, 0x31,
	0x10, 0xfd, 0xd2, 0x8f, 0x56, 0x3a, 0x5a, 0xa5, 0xd1, 0xc3, 0xd2, 0xd3, 0xb2, 0x07, 0x29, 0x88,
	0x29, 0xac, 0x1e, 0x04, 0x4f, 0xb6, 0x0b, 0x3d, 0xaa, 0xf1, 0xe6, 0x2d, 0x9b, 0x1d, 0x6c, 0xb0,
	0x6c, 0xd2, 0x6c, 0xb6, 0xd0, 0x7f, 0xe0, 0xcf, 0xf0, 0xb7, 0xf9, 0x4b, 0xa4, 0xc9, 0xee, 0x41,
	0xca, 0x1e, 0x04, 0x4f, 0xc3, 0x9b, 0x99, 0xf7, 0xf2, 0x5e, 0x18, 0x38, 0xb7, 0x68, 0xb4, 0x75,
	0xb3, 0x50, 0x98, 0xb1, 0xda, 0x69, 0x3a, 0x08, 0x68, 0x32, 0x12, 0x46, 0xcd, 0x84, 0x51, 0xa1,
	0x9d, 0x64, 0x40, 0x9f, 0x6b, 0xb4, 0x3b, 0xee, 0xa7, 0x1c, 0x37, 0x35, 0x56, 0x8e, 0x5e, 0x40,
	0x7f, 0xb3, 0xef, 0x46, 0x24, 0x26, 0xd3, 0x13, 0x1e, 0x00, 0x8d, 0xe0, 0x48, 0x58, 0xb9, 0x52,
	0x5b, 0x8c, 0x7a, 0x31, 0x99, 0x0e, 0x79, 0x0b, 0x93, 0x0c, 0xa2, 0x25, 0xba, 0x47, 0x9b, 0xe1,
	0x1a, 0x1d, 0xfe, 0xd4, 0x3a, 0x85, 0x9e, 0x2a, 0xbc, 0xd0, 0x90, 0xf7, 0x54, 0xd1, 0xad, 0x92,
	0x7e, 0x11, 0x18, 0x05, 0xee, 0x0b, 0xda, 0xad, 0x92, 0x48, 0x6f, 0xa1, 0xef, 0xdd, 0xd1, 0x09,
	0x6b, 0xc2, 0x1c, 0x9a, 0x9d, 0x8c, 0xd8, 0x3e, 0x0e, 0xc7, 0xca, 0xe8, 0xb2, 0xc2, 0xe4, 0x1f,
	0xbd, 0x83, 0xff, 0x4b, 0x74, 0x34, 0x6e, 0x39, 0x5d, 0xd6, 0x0e, 0x99, 0xd7, 0x30, 0x58, 0x58,
	0x14, 0x0e, 0x29, 0xf5, 0xa3, 0x00, 0x3a, 0xd7, 0x53, 0x38, 0x9e, 0x0b, 0x27, 0x57, 0xbf, 0xe0,
	0xa4, 0x0e, 0xc6, 0xc1, 0xc4, 0x42, 0x9b, 0x5d, 0x9b, 0xf3, 0xaa, 0xcd, 0x39, 0xf6, 0xeb, 0x4d,
	0xc8, 0xbf, 0x7b, 0x75, 0xfe, 0x00, 0x97, 0xb2, 0x64, 0x42, 0x59, 0xa5, 0x1d, 0xab, 0x8a, 0x77,
	0x26, 0xd7, 0x0a, 0x4b, 0xc7, 0x8a, 0x3a, 0xcf, 0x35, 0x7b, 0xb3, 0x46, 0x36, 0x9f, 0xf5, 0x44,
	0x5e, 0xcf, 0x58, 0x73, 0x39, 0xf7, 0xa1, 0x7c, 0x10, 0xf2, 0x49, 0x48, 0x3e, 0xf0, 0x07, 0x73,
	0xf3, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x70, 0x14, 0x74, 0x88, 0x5e, 0x02, 0x00, 0x00,
}