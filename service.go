package nrpc

import (
	"fmt"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

func New[T1 any, T2 any](conn *nats.Conn) *service[T1, T2] {
	return &service[T1, T2]{
		conn:             conn,
		requestValidator: func(obj *T1) error { return nil },
	}
}

type service[T1 any, T2 any] struct {
	conn             *nats.Conn
	requestValidator func(obj *T1) error
}

func (m *service[T1, T2]) Handle(msg *nats.Msg, cb func(req *T1) (*T2, error)) (ret *nats.Msg, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("panic:[%v]", err2)
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			log.Printf("\x1b[33mpanic stack info %s\n\x1b[0m", stackInfo)
		}
	}()
	if req, err := decode[T1](msg); err != nil {
		return nil, fmt.Errorf("decode:[%v]", err)
	} else if err := m.requestValidator(req); err != nil {
		return nil, err
	} else if res, err := cb(req); err != nil {
		return nil, err
	} else {
		return encode(res)
	}
}

func (m *service[T1, T2]) Validator(fc func(obj *T1) error) {
	m.requestValidator = fc
}

func (m *service[T1, T2]) Queue(subj string, queue string, cb func(req *T1) (*T2, error)) (*nats.Subscription, error) {
	return m.conn.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		res, err := m.Handle(msg, cb)
		if err != nil {
			res = encodeError(err)
		}
		msg.RespondMsg(res)
	})
}

func (m *service[T1, T2]) Sub(subj string, cb func(req *T1) (*T2, error)) (*nats.Subscription, error) {
	return m.conn.Subscribe(subj, func(msg *nats.Msg) {
		if res, err := m.Handle(msg, cb); err != nil {
			log.Println("handle", err)
			msg.RespondMsg(encodeError(err))
		} else {
			msg.RespondMsg(res)
		}
	})
}
