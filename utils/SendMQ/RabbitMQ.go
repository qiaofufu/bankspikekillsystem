package SendMQ

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	Connection     *amqp.Connection
	ConnectChannel *amqp.Channel
	MessageQueue   amqp.Queue
	QueueNumber    int
}

var MQ RabbitMQ

const (
	OrderQueue = "spike_order"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

func (receiver RabbitMQ) SendMessage(queueName string, messageBody []byte) (err error) {
	if receiver.ConnectChannel == nil {
		receiver.InitRabbitMQ()
	}
	fmt.Println("Send Message")
	err = receiver.ConnectChannel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         messageBody,
			DeliveryMode: amqp.Persistent, //Msg set as persistent
		})
	return
}

func (receiver RabbitMQ) InitRabbitMQ() {
	fmt.Println("Init RabbitMQ")
	username := viper.GetString("rabbitmq.username")
	password := viper.GetString("rabbitmq.password")
	host := viper.GetString("rabbitmq.host")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", username, password, host))
	failOnError(err, "Failed to connect to RabbitMQ")
	MQ.Connection = conn

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	MQ.ConnectChannel = ch

	q, err := ch.QueueDeclare(
		OrderQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	MQ.MessageQueue = q

	MQ.QueueNumber = 1

}

func (receiver RabbitMQ) Close() {
	receiver.Connection.Close()
	receiver.ConnectChannel.Close()
}
