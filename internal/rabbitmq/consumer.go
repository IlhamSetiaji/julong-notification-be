package rabbitmq

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/IlhamSetiaji/julong-notification-be/config"
	"github.com/IlhamSetiaji/julong-notification-be/internal/request"
	"github.com/IlhamSetiaji/julong-notification-be/internal/response"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/IlhamSetiaji/julong-notification-be/utils"
	"github.com/rabbitmq/amqp091-go"
)

func InitConsumer(conf config.Config, log logger.Logger) {
	// conn
	conn, err := amqp091.Dial(conf.RabbitMq.Host)
	if err != nil {
		log.GetLogger().Printf("ERROR: fail init consumer: %s", err.Error())
		os.Exit(1)
	}

	log.GetLogger().Printf("INFO: done init consumer conn")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.GetLogger().Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// create queue
	queue, err := amqpChannel.QueueDeclare(
		conf.RabbitMq.Queue, // channelname
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.GetLogger().Printf("ERROR: fail create queue: %s", err.Error())
		os.Exit(1)
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.GetLogger().Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// consume
	for {
		select {
		case msg := <-msgChannel:
			// unmarshal
			docRply := &response.RabbitMQResponse{}
			docMsg := &request.RabbitMQRequest{}
			err = json.Unmarshal(msg.Body, docRply)
			if err != nil {
				log.GetLogger().Printf("ERROR: fail unmarshl: %s", msg.Body)
				msg.Nack(false, true)
				continue
			}
			log.GetLogger().Printf("INFO: received docRply: %v", docRply)

			err = json.Unmarshal(msg.Body, docMsg)
			if err != nil {
				log.GetLogger().Printf("ERROR: fail unmarshl: %s", msg.Body)
				msg.Nack(false, true)
				continue
			}
			log.GetLogger().Printf("INFO: received docMsg: %v", docMsg)

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.GetLogger().Printf("ERROR: fail to ack: %s", err.Error())
			}

			// find waiting channel(with uid) and forward the reply to it
			if rchan, ok := utils.Rchans[docRply.ID]; ok {
				rchan <- *docRply
			}

			handleMsg(docMsg, log, conf)
		}
	}
}

func handleMsg(docMsg *request.RabbitMQRequest, log logger.Logger, conf config.Config) {
	// switch case
	var msgData map[string]interface{}

	switch docMsg.MessageType {
	case "reply":
		log.GetLogger().Printf("INFO: received reply message")
		return
	default:
		log.GetLogger().Printf("Unknown message type, please recheck your type: %s", docMsg.MessageType)

		msgData = map[string]interface{}{
			"error": errors.New("unknown message type").Error(),
		}
	}
	// reply
	reply := response.RabbitMQResponse{
		ID: docMsg.ID,
		// MessageType: docMsg.MessageType,
		MessageType: "reply",
		MessageData: msgData,
	}
	msg := utils.RabbitMsgConsumer{
		QueueName: docMsg.ReplyTo,
		Reply:     reply,
	}
	utils.Rchan <- msg
}
