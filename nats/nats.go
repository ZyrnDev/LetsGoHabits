package nats

import "github.com/nats-io/nats.go"

type NatsConn = *nats.Conn
type NatsMsg = *nats.Msg

func NatsConnection(natsConnStr string) NatsConn {
	nc, err := nats.Connect(natsConnStr)

	if err != nil {
		panic(err)
	}

	return nc
}
