package nats

import "github.com/nats-io/nats.go"

type NatsConn = *nats.Conn
type NatsMsg = *nats.Msg

type Connection struct {
	conn             *nats.Conn
	connectionString string
}

var connections map[string]Connection = make(map[string]Connection)

func Connect(natsConnStr string) (*Connection, error) {

	connection, connectionExists := connections[natsConnStr]
	if connectionExists {
		return &connection, nil
	}

	nc, err := nats.Connect(natsConnStr)
	if err != nil {
		return nil, err
	}

	connection.connectionString = natsConnStr
	connection.conn = nc

	connections[natsConnStr] = connection

	return &connection, nil
}

func (connection *Connection) Close() {
	connection.conn.Close()
	delete(connections, connection.connectionString)
}

func (connection *Connection) Publish(subject string, data []byte) error {
	return connection.conn.Publish(subject, data)
}

func (connection *Connection) Subscribe(subject string, callback func(msg NatsMsg)) (*nats.Subscription, error) {
	return connection.conn.Subscribe(subject, callback)
}
