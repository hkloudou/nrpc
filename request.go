package nrpc

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type Request[T any] struct {
	Raw  *nats.Msg
	Data *T
}

func NewRequest[T1 any, T2 any](conn *nats.Conn, topic string) *Requester[T1, T2] {
	return &Requester[T1, T2]{
		conn:              conn,
		topic:             topic,
		responseValidator: func(out *Response[T2]) error { return nil },
	}
}

type Requester[T1 any, T2 any] struct {
	conn              *nats.Conn
	topic             string
	responseValidator func(out *Response[T2]) error
}

func (m *Requester[T1, T2]) Validator(fc func(out *Response[T2]) error) *Requester[T1, T2] {
	m.responseValidator = fc
	return m
}

func (m *Requester[T1, T2]) Request(in *T1, timeout time.Duration) (*Response[T2], error) {
	mr, err := encode(in)
	if err != nil {
		return nil, err
	}
	mr.Subject = fmt.Sprintf("%s%s", pre, m.topic)
	res, err := m.conn.RequestMsg(mr, timeout)
	if err != nil {
		return nil, err
	}

	tmp, err := decode[T2](res)
	if err != nil {
		return nil, err
	}
	out := &Response[T2]{
		Raw:  res,
		Data: tmp,
	}
	if err := m.responseValidator(out); err != nil {
		return nil, err
	} else {
		return out, nil
	}
}
