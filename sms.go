package main

import (
	"github.com/spf13/viper"
	"github.com/subosito/twilio"
)

type SMSMessage struct {
	Message		string  `json:"message"`
	To		string  `json:"to"`
	From            string 	`json:"from"`

}

func HandleSMSNotification(smsMessage SMSMessage)  {



	// Initialize twilio Client
	c := twilio.NewClient(viper.GetString(TwilioAccountSid), viper.GetString(TwilioAuthToken), nil)

	// Send Message
	params := twilio.MessageParams{
		Body: smsMessage.Message,
	}

	s, resp, err := c.Messages.Send("6783039533", smsMessage.To, params)
	log.Info("Send:", s)
	log.Info("Response:", resp)
	log.Info("Err:", err)
	
}
