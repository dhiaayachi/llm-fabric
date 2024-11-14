package llm

import (
	"context"
	"encoding/json"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Llm interface {
	// SubmitTask TODO: task need to be changed to proto eventually
	SubmitTask(ctx context.Context, task string, opts ...*agentinfo.LlmOpt) (response string, err error)
}

func getOpt[T any](typ agentinfo.LlmOptType, opts ...*agentinfo.LlmOpt) T {
	var empty T
	for _, o := range opts {
		if o.Typ == typ {

			val, err := GetVal[T](o)
			if err != nil {
				return empty
			}
			return val
		}
	}
	return empty
}

func GetVal[T any](o *agentinfo.LlmOpt) (T, error) {
	var value T
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

func FromVal[T any](o *agentinfo.LlmOpt, v T) error {
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
