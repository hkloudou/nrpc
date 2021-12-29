package nrpc

import (
	"time"

	"github.com/nats-io/nats.go"
)

func NewRequest[T1 any, T2 any](conn *nats.Conn, topic string) *request[T1, T2] {
	return &request[T1, T2]{
		conn:              conn,
		topic:             topic,
		responseValidator: func(obj *T2) error { return nil },
	}
}

type request[T1 any, T2 any] struct {
	conn              *nats.Conn
	topic             string
	responseValidator func(obj *T2) error
}

func (m *request[T1, T2]) Validator(fc func(obj *T2) error) {
	m.responseValidator = fc
}

func (m *request[T1, T2]) Request(in *T1, timeout time.Duration) (*T2, error) {
	mr, err := encode(in)
	if err != nil {
		return nil, err
	}
	mr.Subject = m.topic
	res, err := m.conn.RequestMsg(mr, timeout)
	if err != nil {
		return nil, err
	}

	if tmp, err := decode[T2](res); err != nil {
		return nil, err
	} else if err := m.responseValidator(tmp); err != nil {
		return nil, err
	} else {
		return tmp, nil
	}
}
