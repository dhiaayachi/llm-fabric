// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: agent/v1/agent.proto

package agentv1

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

type Capability int32

const (
	Capability_CAPABILITY_UNSPECIFIED    Capability = 0
	Capability_CAPABILITY_TEXT           Capability = 1
	Capability_CAPABILITY_SUMMARIZATION  Capability = 2
	Capability_CAPABILITY_IMAGE          Capability = 3
	Capability_CAPABILITY_VIDEO          Capability = 4
	Capability_CAPABILITY_DISPATCH       Capability = 5
	Capability_CAPABILITY_SUBDIVIDE_TASK Capability = 6
)

// Enum value maps for Capability.
var (
	Capability_name = map[int32]string{
		0: "CAPABILITY_UNSPECIFIED",
		1: "CAPABILITY_TEXT",
		2: "CAPABILITY_SUMMARIZATION",
		3: "CAPABILITY_IMAGE",
		4: "CAPABILITY_VIDEO",
		5: "CAPABILITY_DISPATCH",
		6: "CAPABILITY_SUBDIVIDE_TASK",
	}
	Capability_value = map[string]int32{
		"CAPABILITY_UNSPECIFIED":    0,
		"CAPABILITY_TEXT":           1,
		"CAPABILITY_SUMMARIZATION":  2,
		"CAPABILITY_IMAGE":          3,
		"CAPABILITY_VIDEO":          4,
		"CAPABILITY_DISPATCH":       5,
		"CAPABILITY_SUBDIVIDE_TASK": 6,
	}
)

func (x Capability) Enum() *Capability {
	p := new(Capability)
	*p = x
	return p
}

func (x Capability) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Capability) Descriptor() protoreflect.EnumDescriptor {
	return file_agent_v1_agent_proto_enumTypes[0].Descriptor()
}

func (Capability) Type() protoreflect.EnumType {
	return &file_agent_v1_agent_proto_enumTypes[0]
}

func (x Capability) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Capability.Descriptor instead.
func (Capability) EnumDescriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{0}
}

type ToolType int32

const (
	ToolType_TOOL_TYPE_UNSPECIFIED ToolType = 0
	ToolType_TOOL_TYPE_RAG         ToolType = 1
	ToolType_TOOL_TYPE_FUNCTION    ToolType = 2
)

// Enum value maps for ToolType.
var (
	ToolType_name = map[int32]string{
		0: "TOOL_TYPE_UNSPECIFIED",
		1: "TOOL_TYPE_RAG",
		2: "TOOL_TYPE_FUNCTION",
	}
	ToolType_value = map[string]int32{
		"TOOL_TYPE_UNSPECIFIED": 0,
		"TOOL_TYPE_RAG":         1,
		"TOOL_TYPE_FUNCTION":    2,
	}
)

func (x ToolType) Enum() *ToolType {
	p := new(ToolType)
	*p = x
	return p
}

func (x ToolType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ToolType) Descriptor() protoreflect.EnumDescriptor {
	return file_agent_v1_agent_proto_enumTypes[1].Descriptor()
}

func (ToolType) Type() protoreflect.EnumType {
	return &file_agent_v1_agent_proto_enumTypes[1]
}

func (x ToolType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ToolType.Descriptor instead.
func (ToolType) EnumDescriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{1}
}

type Agent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// port is the port assigned to the workload
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// description opaque description
	Description  string       `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Capabilities []Capability `protobuf:"varint,4,rep,packed,name=capabilities,proto3,enum=agent.v1.Capability" json:"capabilities,omitempty"`
	Tools        []*Tool      `protobuf:"bytes,5,rep,name=tools,proto3" json:"tools,omitempty"`
	Id           string       `protobuf:"bytes,6,opt,name=id,proto3" json:"id,omitempty"`
	Cost         float32      `protobuf:"fixed32,7,opt,name=cost,proto3" json:"cost,omitempty"`
	Score        float32      `protobuf:"fixed32,8,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *Agent) Reset() {
	*x = Agent{}
	mi := &file_agent_v1_agent_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Agent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Agent) ProtoMessage() {}

func (x *Agent) ProtoReflect() protoreflect.Message {
	mi := &file_agent_v1_agent_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Agent.ProtoReflect.Descriptor instead.
func (*Agent) Descriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{0}
}

func (x *Agent) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Agent) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Agent) GetCapabilities() []Capability {
	if x != nil {
		return x.Capabilities
	}
	return nil
}

func (x *Agent) GetTools() []*Tool {
	if x != nil {
		return x.Tools
	}
	return nil
}

func (x *Agent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Agent) GetCost() float32 {
	if x != nil {
		return x.Cost
	}
	return 0
}

func (x *Agent) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

type Tool struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Type        ToolType `protobuf:"varint,3,opt,name=type,proto3,enum=agent.v1.ToolType" json:"type,omitempty"`
}

func (x *Tool) Reset() {
	*x = Tool{}
	mi := &file_agent_v1_agent_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tool) ProtoMessage() {}

func (x *Tool) ProtoReflect() protoreflect.Message {
	mi := &file_agent_v1_agent_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tool.ProtoReflect.Descriptor instead.
func (*Tool) Descriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{1}
}

func (x *Tool) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tool) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Tool) GetType() ToolType {
	if x != nil {
		return x.Type
	}
	return ToolType_TOOL_TYPE_UNSPECIFIED
}

type SubmitTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Task string `protobuf:"bytes,1,opt,name=task,proto3" json:"task,omitempty"`
}

func (x *SubmitTaskRequest) Reset() {
	*x = SubmitTaskRequest{}
	mi := &file_agent_v1_agent_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitTaskRequest) ProtoMessage() {}

func (x *SubmitTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_v1_agent_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitTaskRequest.ProtoReflect.Descriptor instead.
func (*SubmitTaskRequest) Descriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{2}
}

func (x *SubmitTaskRequest) GetTask() string {
	if x != nil {
		return x.Task
	}
	return ""
}

type SubmitTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *SubmitTaskResponse) Reset() {
	*x = SubmitTaskResponse{}
	mi := &file_agent_v1_agent_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitTaskResponse) ProtoMessage() {}

func (x *SubmitTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_v1_agent_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitTaskResponse.ProtoReflect.Descriptor instead.
func (*SubmitTaskResponse) Descriptor() ([]byte, []int) {
	return file_agent_v1_agent_proto_rawDescGZIP(), []int{3}
}

func (x *SubmitTaskResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

var File_agent_v1_agent_proto protoreflect.FileDescriptor

var file_agent_v1_agent_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x22, 0xdd, 0x01, 0x0a, 0x05, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x0c, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x52, 0x0c, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73,
	0x12, 0x24, 0x0a, 0x05, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6f, 0x6c, 0x52,
	0x05, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63,
	0x6f, 0x72, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x22, 0x64, 0x0a, 0x04, 0x54, 0x6f, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x27, 0x0a, 0x11, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x61, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x22,
	0x30, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2a, 0xbf, 0x01, 0x0a, 0x0a, 0x43, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x12, 0x1a, 0x0a, 0x16, 0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f,
	0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x10,
	0x01, 0x12, 0x1c, 0x0a, 0x18, 0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f,
	0x53, 0x55, 0x4d, 0x4d, 0x41, 0x52, 0x49, 0x5a, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x12,
	0x14, 0x0a, 0x10, 0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x49, 0x4d,
	0x41, 0x47, 0x45, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c,
	0x49, 0x54, 0x59, 0x5f, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x10, 0x04, 0x12, 0x17, 0x0a, 0x13, 0x43,
	0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54,
	0x43, 0x48, 0x10, 0x05, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49,
	0x54, 0x59, 0x5f, 0x53, 0x55, 0x42, 0x44, 0x49, 0x56, 0x49, 0x44, 0x45, 0x5f, 0x54, 0x41, 0x53,
	0x4b, 0x10, 0x06, 0x2a, 0x50, 0x0a, 0x08, 0x54, 0x6f, 0x6f, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x19, 0x0a, 0x15, 0x54, 0x4f, 0x4f, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x54, 0x4f,
	0x4f, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x41, 0x47, 0x10, 0x01, 0x12, 0x16, 0x0a,
	0x12, 0x54, 0x4f, 0x4f, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x46, 0x55, 0x4e, 0x43, 0x54,
	0x49, 0x4f, 0x4e, 0x10, 0x02, 0x32, 0x59, 0x0a, 0x0c, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x0a, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x54,
	0x61, 0x73, 0x6b, 0x12, 0x1b, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d,
	0x69, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x98, 0x01, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x42, 0x0a, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x68, 0x69, 0x61,
	0x61, 0x79, 0x61, 0x63, 0x68, 0x69, 0x2f, 0x6c, 0x6c, 0x6d, 0x2d, 0x66, 0x61, 0x62, 0x72, 0x69,
	0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41,
	0x58, 0x58, 0xaa, 0x02, 0x08, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x08,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x14, 0x41, 0x67, 0x65, 0x6e, 0x74,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x09, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_agent_v1_agent_proto_rawDescOnce sync.Once
	file_agent_v1_agent_proto_rawDescData = file_agent_v1_agent_proto_rawDesc
)

func file_agent_v1_agent_proto_rawDescGZIP() []byte {
	file_agent_v1_agent_proto_rawDescOnce.Do(func() {
		file_agent_v1_agent_proto_rawDescData = protoimpl.X.CompressGZIP(file_agent_v1_agent_proto_rawDescData)
	})
	return file_agent_v1_agent_proto_rawDescData
}

var file_agent_v1_agent_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_agent_v1_agent_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_agent_v1_agent_proto_goTypes = []any{
	(Capability)(0),            // 0: agent.v1.Capability
	(ToolType)(0),              // 1: agent.v1.ToolType
	(*Agent)(nil),              // 2: agent.v1.Agent
	(*Tool)(nil),               // 3: agent.v1.Tool
	(*SubmitTaskRequest)(nil),  // 4: agent.v1.SubmitTaskRequest
	(*SubmitTaskResponse)(nil), // 5: agent.v1.SubmitTaskResponse
}
var file_agent_v1_agent_proto_depIdxs = []int32{
	0, // 0: agent.v1.Agent.capabilities:type_name -> agent.v1.Capability
	3, // 1: agent.v1.Agent.tools:type_name -> agent.v1.Tool
	1, // 2: agent.v1.Tool.type:type_name -> agent.v1.ToolType
	4, // 3: agent.v1.AgentService.SubmitTask:input_type -> agent.v1.SubmitTaskRequest
	5, // 4: agent.v1.AgentService.SubmitTask:output_type -> agent.v1.SubmitTaskResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_agent_v1_agent_proto_init() }
func file_agent_v1_agent_proto_init() {
	if File_agent_v1_agent_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_agent_v1_agent_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_agent_v1_agent_proto_goTypes,
		DependencyIndexes: file_agent_v1_agent_proto_depIdxs,
		EnumInfos:         file_agent_v1_agent_proto_enumTypes,
		MessageInfos:      file_agent_v1_agent_proto_msgTypes,
	}.Build()
	File_agent_v1_agent_proto = out.File
	file_agent_v1_agent_proto_rawDesc = nil
	file_agent_v1_agent_proto_goTypes = nil
	file_agent_v1_agent_proto_depIdxs = nil
}
