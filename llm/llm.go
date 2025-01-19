package llm

import (
	"context"
	"encoding/json"
	llmoptions "github.com/dhiaayachi/llm-fabric/proto/gen/llm_options/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Llm interface {
	// SubmitTask TODO: task need to be changed to proto eventually
	SubmitTask(ctx context.Context, task string, opts ...*llmoptions.LlmOpt) (response string, err error)
	SubmitTaskWithSchema(ctx context.Context, task string, schema string) (response string, err error)
}

func getOpt[T any](typ llmoptions.LlmOptType, opts ...*llmoptions.LlmOpt) T {
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

func GetVal[T any](o *llmoptions.LlmOpt) (T, error) {
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

func FromVal[T any](o *llmoptions.LlmOpt, v T) error {
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
