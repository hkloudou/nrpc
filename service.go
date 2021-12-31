package nrpc

import (
	"fmt"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

const pre = "_NRPC."

func New[T1 any, T2 any](conn *nats.Conn) *Servicer[T1, T2] {
	return &Servicer[T1, T2]{
		conn:             conn,
		requestValidator: func(req *Request[T1]) error { return nil },
	}
}

type Response[T any] struct {
	Raw  *nats.Msg
	Data *T
}

type Servicer[T1 any, T2 any] struct {
	conn             *nats.Conn
	requestValidator func(in *Request[T1]) error
}

func (m *Servicer[T1, T2]) Handle(msg *nats.Msg, cb func(in *Request[T1]) (*T2, error)) (ret *nats.Msg, err error) {
	con, _ := m.conn.JetStream()
	con.AddStream(&nats.StreamConfig{})

	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("nats: panic[%v]", err2)
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			log.Printf("\x1b[33mpanic stack info %s\n\x1b[0m", stackInfo)
		}
	}()
	tmp, err := decode[T1](msg)
	if err != nil {
		return nil, fmt.Errorf("decode:[%v]", err)
	}
	req := &Request[T1]{
		Raw:  msg,
		Data: tmp,
	}
	if err := m.requestValidator(req); err != nil {
		return nil, err
	} else if res, err := cb(req); err != nil {
		return nil, err
	} else {
		return encode(res)
	}
}

func (m *Servicer[T1, T2]) Validator(fc func(in *Request[T1]) error) *Servicer[T1, T2] {
	m.requestValidator = fc
	return m
}

func (m *Servicer[T1, T2]) Queue(subj string, queue string, cb func(in *Request[T1]) (*T2, error)) (*nats.Subscription, error) {
	return m.conn.QueueSubscribe(fmt.Sprintf("%s%s", pre, subj), queue, func(msg *nats.Msg) {
		res, err := m.Handle(msg, cb)
		if err != nil {
			res = encodeError(err)
		}
		msg.RespondMsg(res)
	})
}

func (m *Servicer[T1, T2]) Sub(subj string, cb func(in *Request[T1]) (*T2, error)) (*nats.Subscription, error) {
	return m.conn.Subscribe(fmt.Sprintf("%s%s", pre, subj), func(msg *nats.Msg) {
		if res, err := m.Handle(msg, cb); err != nil {
			log.Println("handle", err)
			msg.RespondMsg(encodeError(err))
		} else {
			msg.RespondMsg(res)
		}
	})
}
