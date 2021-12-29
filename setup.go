package nrpc

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

var _conn *nats.Conn

func Connect(url string, options ...nats.Option) (*nats.Conn, error) {
	if url == "" {
		url = nats.DefaultURL
	}
	if len(options) == 0 {
		options = DefaultConfig()
	}
	var err error
	_conn, err = nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	return _conn, nil
}
func GetConn() *nats.Conn {
	return _conn
}

func DefaultConfig() []nats.Option {
	return []nats.Option{
		nats.PingInterval(5 * time.Second),
		nats.MaxPingsOutstanding(3),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Println("Got disconnected! Reason: ", err)
		}),
		nats.MaxReconnects(60),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			if nc.Reconnects > 58 {
				panic("too much times")
			}
			log.Println("Got reconnected to", "[", nc.Reconnects, "]", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Println("Connection closed. Reason: ", nc.LastError())
		}),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			log.Println("Discover closed. Reason: ", nc.ConnectedAddr(), nc.ConnectedUrl())
		}),
		nats.CustomReconnectDelay(func(attempts int) time.Duration {
			return time.Second * 1
		}),
	}
}
