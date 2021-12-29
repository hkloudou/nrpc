package nrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func PointerOf[T any](v T) *T {
	return &v
}

func decode[T any](msg *nats.Msg) (*T, error) {
	if msg.Header.Get("Nil") == "1" {
		return nil, nil
	}
	var data T
	mt := reflect.ValueOf(&data).MethodByName("ProtoReflect")
	headErr := msg.Header.Get("Error")
	if len(headErr) != 0 {
		return nil, errors.New(headErr)
	}
	if mt.IsValid() {
		if err := proto.Unmarshal(msg.Data, protoreflect.ValueOf(reflect.ValueOf(&data).MethodByName("ProtoReflect").Call(nil)[0].Interface()).Message().Interface()); err != nil {
			return nil, err
		}
		return &data, nil
	} else {
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return nil, err
		}
		return &data, nil
	}
}

func encode[T any](obj *T) (*nats.Msg, error) {
	if obj == nil {
		return &nats.Msg{
			Header: nats.Header{
				"Nil": []string{"1"},
			},
		}, nil
	}
	var b []byte
	var err = fmt.Errorf("not support")
	mt :=
		reflect.ValueOf(obj).MethodByName("ProtoReflect")
	if mt.IsValid() {
		b, err = proto.Marshal(protoreflect.ValueOf(mt.Call(nil)[0].Interface()).Message().Interface())
	} else {
		b, err = json.Marshal(obj)
	}
	if err != nil {
		return nil, err
	}
	return &nats.Msg{
		Header: make(nats.Header),
		Data:   b,
	}, nil
}

func encodeError(err error) *nats.Msg {
	if err == nil {
		return &nats.Msg{Header: nats.Header{"Error": []string{"error"}}}
	}
	return &nats.Msg{Header: nats.Header{"Error": []string{err.Error()}}}
}
