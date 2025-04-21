package bus

import (
	"github.com/nats-io/nats.go"
)

type natsBroker struct {
	client *nats.Conn
}

func NewNatsBroker(url string) (Driver, error) {
	c, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return natsBroker{c}, nil
}

func (n natsBroker) Publish(subject string, data []byte) error {
	return n.client.Publish(subject, data)
}

func (n natsBroker) Subscribe(subject string) (<-chan []byte, <-chan error, error) {
	dataCh := make(chan []byte)
	errCh := make(chan error)
	_, err := n.client.Subscribe(subject, func(msg *nats.Msg) {
		if err := msg.Ack(); err != nil {
			errCh <- err
		}
		dataCh <- msg.Data
	})
	n.client.SetClosedHandler(func(c *nats.Conn) {
		close(dataCh)
		close(errCh)
	})
	return dataCh, errCh, err
}

func (n natsBroker) Close() error {
	n.client.Close()
	return nil
}
