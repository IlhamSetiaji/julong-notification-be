package rabbitmq

import (
	"encoding/json"
	"os"

	"github.com/IlhamSetiaji/julong-notification-be/config"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/IlhamSetiaji/julong-notification-be/utils"
	"github.com/rabbitmq/amqp091-go"
)

func InitProducer(conf config.Config, log logger.Logger) {
	// conn
	conn, err := amqp091.Dial(conf.RabbitMq.Host)
	if err != nil {
		log.GetLogger().Printf("ERROR: fail init consumer: %s", err.Error())
		os.Exit(1)
	}

	log.GetLogger().Printf("INFO: done init producer conn")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.GetLogger().Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	for {
		select {
		case msg := <-utils.Pchan:
			// marshal
			data, err := json.Marshal(&msg.Message)
			if err != nil {
				log.GetLogger().Printf("ERROR: fail marshal: %s", err.Error())
				continue
			}

			// publish message
			err = amqpChannel.Publish(
				"",            // exchange
				msg.QueueName, // routing key
				false,         // mandatory
				false,         // immediate
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        data,
				},
			)
			if err != nil {
				log.GetLogger().Printf("ERROR: fail publish msg: %s", err.Error())
				continue
			}

			log.GetLogger().Printf("INFO: published msg: %v", msg.Message)
		case msg := <-utils.Rchan:
			// marshal
			data, err := json.Marshal(&msg.Reply)
			if err != nil {
				log.GetLogger().Printf("ERROR: fail marshal: %s", err.Error())
				continue
			}

			// publish message
			err = amqpChannel.Publish(
				"",            // exchange
				msg.QueueName, // routing key
				false,         // mandatory
				false,         // immediate
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        data,
				},
			)
			if err != nil {
				log.GetLogger().Printf("ERROR: fail publish msg: %s", err.Error())
				continue
			}

			log.GetLogger().Printf("INFO: published msg: %v to: %s", msg.Reply, msg.QueueName)
		}
	}
}
