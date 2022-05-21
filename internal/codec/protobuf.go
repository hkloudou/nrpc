package codec

import (
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type protobufCodec struct {
}

func init() {
	RegisteCodec("protobuf", &jsonCodec{})
}

func (m *protobufCodec) Check(val interface{}) bool {
	return reflect.ValueOf(val).MethodByName("ProtoReflect").IsValid()
}

func (m *protobufCodec) Marshal(val interface{}) ([]byte, error) {
	mt := reflect.ValueOf(val).MethodByName("ProtoReflect")
	return proto.Marshal(protoreflect.ValueOf(mt.Call(nil)[0].Interface()).Message().Interface())
}

func (m *protobufCodec) Unmarshal(data []byte, val interface{}) error {
	return proto.Unmarshal(data, protoreflect.ValueOf(reflect.ValueOf(&data).MethodByName("ProtoReflect").Call(nil)[0].Interface()).Message().Interface())
}
