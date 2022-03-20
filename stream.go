package nrpc

import (
	"errors"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

type Stream[Tr any] struct {
	js nats.JetStreamContext
	// topic             string
	// requestValidator  func(in *Request[T1]) error
	// responseValidator func(out *Response[T2]) error
}

func NewStream[T any](conn *nats.Conn, opts ...nats.JSOpt) (*Stream[T], error) {
	js, err := conn.JetStream(conn, opts)
	if err != nil {
		return nil, err
	}
	return &Stream[T]{js: js}, nil
}

func (m *Stream[T]) Config(cfg *nats.StreamConfig) error {
	if m.js == nil {
		return errors.New("please InitNats")
	}
	stream, err := m.js.StreamInfo(cfg.Name)
	if err != nil && !strings.Contains(err.Error(), "stream not found") {
		return err
	}
	//stream not found
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", cfg.Name, cfg.Subjects)
		_, err = m.js.AddStream(cfg)
		if err != nil {
			return err
		}
	} else {
		_, err = m.js.UpdateStream(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Stream[T]) Js() nats.JetStreamContext {
	return m.js
}
