package bus

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbit struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitBroker(url string, exqs map[string][]string) (Driver, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	for ex, qs := range exqs {
		if err := channel.ExchangeDeclare(ex, "direct", true, false, false, false, nil); err != nil {
			return nil, err
		}
		for _, q := range qs {
			_, err := channel.QueueDeclare(q, true, true, false, false, nil)
			if err != nil {
				return nil, err
			}
			if err = channel.QueueBind(q, "", ex, false, nil); err != nil {
				return nil, err
			}
		}
	}
	return rabbit{conn: conn, ch: channel}, nil
}

func (r rabbit) Publish(subject string, data []byte) error {
	return r.ch.Publish(subject, "", false, false, amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        data,
		Timestamp:   time.Now(),
	})
}

func (r rabbit) Subscribe(queue string) (<-chan []byte, <-chan error, error) {
	data, err := r.ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	dataCh := make(chan []byte)
	go func() {
		defer close(dataCh)
		for d := range data {
			dataCh <- d.Body
		}
	}()
	errCh := make(chan error)
	close(errCh)
	return dataCh, errCh, nil
}

func (r rabbit) Close() error {
	return errors.Join(r.ch.Close(), r.conn.Close())
}
