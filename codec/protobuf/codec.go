package protobuf

import (
	"reflect"

	"github.com/hkloudou/nrpc/face"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type jsonCodec struct {
}

func init() {
	face.RegisteCodec("protobuf", &jsonCodec{})
}

func (m *jsonCodec) Check(val interface{}) bool {
	return reflect.ValueOf(val).MethodByName("ProtoReflect").IsValid()
}

func (m *jsonCodec) Marshal(val interface{}) ([]byte, error) {
	mt := reflect.ValueOf(val).MethodByName("ProtoReflect")
	return proto.Marshal(protoreflect.ValueOf(mt.Call(nil)[0].Interface()).Message().Interface())
}

func (m *jsonCodec) Unmarshal(data []byte, val interface{}) error {
	return proto.Unmarshal(data, protoreflect.ValueOf(reflect.ValueOf(&data).MethodByName("ProtoReflect").Call(nil)[0].Interface()).Message().Interface())
}
