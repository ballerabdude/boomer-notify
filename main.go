package main

import (
	"github.com/spf13/viper"
	"fmt"
	"gopkg.in/op/go-logging.v1"
	"github.com/streadway/amqp"
	"encoding/json"
)
const (

	AmqpConnectionUrl               = "amqp.connection_url"
)

type AmqpMessage struct {

	Type    	string   `json:"type"`
}

func init()  {

	// Configurations are loaded from a config.yaml file which is required but not checked in
	// Get it from  abdul.hagi@turner.com
	// Set up our configurations / Environment var
	log.Info("Loading Configuration")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.BindEnv(AmqpConnectionUrl, "AmqpConnectionUrl")
	viper.AutomaticEnv()


}

var (
	//log
	log = logging.MustGetLogger("artemis")


)

func main() {

	// Connect to Rabbitmq server
	log.Info("Connecting to RabbitMQ")
	conn, err := amqp.Dial(viper.GetString(AmqpConnectionUrl))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"tasks",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Info("Received a message: %s", d.Body)
			var amqpMessage AmqpMessage
			json.Unmarshal(d.Body, &amqpMessage)
			log.Info("Message Type: %s", amqpMessage.Type)
		}
	}()

	log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
