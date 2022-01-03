package nrpc

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

const pre = "_NRPC."

func New[T1 any, T2 any](conn *nats.Conn, topic string) *Servicer[T1, T2] {
	return &Servicer[T1, T2]{
		conn:  conn,
		topic: topic,
	}
}

type Response[T any] struct {
	Raw *nats.Msg

	Data *T
}

type Servicer[T1 any, T2 any] struct {
	conn              *nats.Conn
	topic             string
	requestValidator  func(in *Request[T1]) error
	responseValidator func(out *Response[T2]) error
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
	if m.requestValidator != nil {
		if err := m.requestValidator(req); err != nil {
			return nil, err
		}
	}
	if res, err := cb(req); err != nil {
		return nil, err
	} else {
		return encode(res)
	}
}

func (m *Servicer[T1, T2]) Validator(fc func(in *Request[T1]) error) *Servicer[T1, T2] {
	m.requestValidator = fc
	return m
}

func (m *Servicer[T1, T2]) Validator2(fc func(out *Response[T2]) error) *Servicer[T1, T2] {
	m.responseValidator = fc
	return m
}

func (m *Servicer[T1, T2]) Request(in *T1, timeout time.Duration) (*Response[T2], error) {
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
	if m.responseValidator != nil {
		if err := m.responseValidator(out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (m *Servicer[T1, T2]) Queue(cb func(in *Request[T1]) (*T2, error), queues ...string) ([]*nats.Subscription, error) {
	rets := make([]*nats.Subscription, 0)
	arr := queues
	if len(arr) == 0 {
		arr = append(arr, "worker")
	}
	for i := 0; i < len(arr); i++ {
		res, err := m.conn.QueueSubscribe(fmt.Sprintf("%s%s", pre, m.topic), arr[i], func(msg *nats.Msg) {
			res, err := m.Handle(msg, cb)
			if err != nil {
				res = encodeError(err)
			}
			msg.RespondMsg(res)
		})
		if err != nil {
			return nil, err
		}
		rets = append(rets, res)
	}
	return rets, nil
}

func (m *Servicer[T1, T2]) Sub(cb func(in *Request[T1]) (*T2, error)) (*nats.Subscription, error) {
	return m.conn.Subscribe(fmt.Sprintf("%s%s", pre, m.topic), func(msg *nats.Msg) {
		if res, err := m.Handle(msg, cb); err != nil {
			msg.RespondMsg(encodeError(err))
		} else {
			msg.RespondMsg(res)
		}
	})
}
