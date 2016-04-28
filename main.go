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
	SMSChannel			= "amqp.channels.sms"
	EmailQueue			= "amqp.channels.email"

	TwilioAccountSid	 	= "twilio.account_sid"
	TwilioAuthToken			= "twilio.auth_token"
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
	viper.BindEnv(SMSChannel, "SMSChannel")
	viper.BindEnv(EmailQueue, "EmailChannel")

	viper.BindEnv(TwilioAccountSid, "TwilioAccountSid")
	viper.BindEnv(TwilioAuthToken, "TwilioAuthToken")

	viper.AutomaticEnv()


}

var (
	//log
	log = logging.MustGetLogger("artemis")

	// amqp channel
	AmpqChannel *amqp.Channel

	//amqp connection
	AmqpConnection *amqp.Connection


)

func main() {

	// Connect to Rabbitmq server
	log.Info("Connecting to RabbitMQ")
	conn, err := amqp.Dial(viper.GetString(AmqpConnectionUrl))
	AmqpConnection = conn
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	AmpqChannel = ch
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	listenToSMS()
	listenToEmail()
	log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)

	<-forever


}

func listenToSMS()  {
	q, err := AmpqChannel.QueueDeclare(
		viper.GetString(SMSChannel), // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := AmpqChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")


	go func() {
		for d := range msgs {
			log.Info("Received a message: %s", d.Body)
			var smsMessage SMSMessage
			json.Unmarshal(d.Body, &smsMessage)
			log.Info("SMS Message: %s", smsMessage.Message)
			HandleSMSNotification(smsMessage)
		}
	}()



}

func listenToEmail()  {
	q, err := AmpqChannel.QueueDeclare(
		viper.GetString(EmailQueue), // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := AmpqChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")


	go func() {
		for d := range msgs {
			log.Info("Received a message: %s", d.Body)
			var smsMessage SMSMessage
			json.Unmarshal(d.Body, &smsMessage)
			log.Info("Email Message: %s", smsMessage.Message)
		}
	}()

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
