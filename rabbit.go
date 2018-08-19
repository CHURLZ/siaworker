package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func Connect() *Connection {
	user := os.Getenv("RABBIT_USER")
	passwd := os.Getenv("RABBIT_PASSWD")
	url := os.Getenv("RABBIT_URL")
	port := os.Getenv("RABBIT_PORT")

	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, passwd, url, port)
	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return &Connection{conn, ch, q}
}

func (conn *Connection) Publish(msg []byte) error {
	return conn.channel.Publish(
		"",              // exchange
		conn.queue.Name, // routing key
		false,           // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/xml",
			Body:         []byte(msg),
		})
}

func (conn *Connection) Close() error {
	err := conn.conn.Close()
	failOnError(err, "Failed to close connection.")

	err = conn.channel.Close()
	failOnError(err, "Failed to close channel.")

	return err
}
