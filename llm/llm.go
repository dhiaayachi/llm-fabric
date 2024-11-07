package llm

import (
	"context"
	"encoding/json"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Llm interface {
	// SubmitTask TODO: task need to be changed to proto eventually
	SubmitTask(ctx context.Context, task string, opts ...*Opt) (response string, err error)
	GetCapabilities() []agentv1.Capability // Abilities or features the llm supports
	GetTools() []agentv1.Tool              // Abilities or features the llm supports
}

func getOpt[T any](typ agentv1.LlmOptType, opts ...*Opt) T {
	var empty T
	for _, o := range opts {
		if o.Typ == typ {
			val, err := o.GetVal()
			if err != nil {
				return empty
			}
			return val.(T)
		}
	}
	return empty
}

type Opt struct {
	*agentv1.LlmOpt
}

func (o *Opt) GetVal() (interface{}, error) {
	var value interface{}
	bytesValue := &wrapperspb.BytesValue{}
	err := anypb.UnmarshalTo(o.GetLlmOptVal(), bytesValue, proto.UnmarshalOptions{})
	if err != nil {
		return value, err
	}
	uErr := json.Unmarshal(bytesValue.Value, &value)
	if uErr != nil {
		return value, uErr
	}
	return value, nil
}

func (o *Opt) FromVal(v interface{}) error {
	anyValue := &anypb.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrapperspb.BytesValue{
		Value: bytes,
	}
	err := anypb.MarshalFrom(anyValue, bytesValue, proto.MarshalOptions{})
	if err != nil {
		return err
	}
	o.LlmOptVal = anyValue
	return nil
}
