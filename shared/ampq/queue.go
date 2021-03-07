package ampq

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Conn struct {
	Channel *amqp.Channel
}

func GetQueueConnection(rabbitURI string) (Conn, error) {
	conn, err := amqp.Dial(rabbitURI)
	if err != nil {
		log.Println(err)
		return Conn{}, errors.New("unable to contact remote RabbitMQ service")
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return Conn{}, errors.New("unable to build channel information")
	}
	return Conn{
		Channel: channel,
	}, nil
}

func (conn Conn) StartConsumer(queueName string, routingKey string, handler func(d amqp.Delivery) bool, concurrency int) error {
	// Boilerplate pulled from https://qvault.io/2020/04/29/connecting-to-rabbitmq-in-golang/#:~:text=RabbitMQ%20is%20a%20great%20message%20broker%20with%20awesome,social%20media%20posts%20through%20our%20Go%20services%20daily.

	// create the queue if it doesn't already exist
	_, err := conn.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// bind the queue to the routing key
	err = conn.Channel.QueueBind(queueName, routingKey, "events", false, nil)
	if err != nil {
		return err
	}

	// prefetch 4x as many messages as we can handle at once
	prefetchCount := concurrency * 4
	err = conn.Channel.Qos(prefetchCount, 0, false)
	if err != nil {
		return err
	}

	msgs, err := conn.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// create a goroutine for the number of concurrent threads requested
	for i := 0; i < concurrency; i++ {
		log.Printf("Processing messages on thread %v...\n", i)
		go func() {
			for msg := range msgs {
				// if tha handler returns true then ACK, else NACK
				// the message back into the rabbit queue for
				// another round of processing
				if handler(msg) {
					msg.Ack(false)
				} else {
					msg.Nack(false, true)
				}
			}
			log.Println("Rabbit consumer closed - critical Error")
			os.Exit(1)
		}()
	}
	return nil
}
