package nrpc

import (
	"errors"
	"fmt"

	"github.com/hkloudou/nrpc/internal/codec"
	"github.com/nats-io/nats.go"
)

func PointerOf[T any](v T) *T {
	return &v
}

func decode[T any](msg *nats.Msg) (*T, error) {
	if msg.Header.Get("Nil") == "1" {
		return nil, nil
	}
	var data T
	headErr := msg.Header.Get("Error")
	if len(headErr) != 0 {
		return nil, errors.New(headErr)
	}
	code := codec.GetCodec(&data)
	if err := code.Unmarshal(msg.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
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

	b, err = codec.GetCodec(obj).Marshal(obj)
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
