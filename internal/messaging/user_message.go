package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-notification-be/internal/request"
	"github.com/IlhamSetiaji/julong-notification-be/internal/response"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/IlhamSetiaji/julong-notification-be/utils"
	"github.com/google/uuid"
)

type IUserMessage interface {
	SendFindUserByIDMessage(request request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error)
	SendGetUserMe(request request.SendFindUserByIDMessageRequest) (*response.SendGetUserMeResponse, error)
}

type UserMessage struct {
	Log logger.Logger
}

func NewUserMessage(log logger.Logger) IUserMessage {
	return &UserMessage{
		Log: log,
	}
}

func UserMessageFactory(log logger.Logger) IUserMessage {
	return NewUserMessage(log)
}

func (m *UserMessage) SendFindUserByIDMessage(req request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_user_by_id",
		MessageData: payload,
		ReplyTo:     "julong_manpower",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsgPublisher{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindUserByIDMessage] " + errMsg)
	}

	return &response.SendFindUserByIDResponse{
		ID:   resp.MessageData["user_id"].(string),
		Name: resp.MessageData["name"].(string),
	}, nil
}

func (m *UserMessage) SendGetUserMe(req request.SendFindUserByIDMessageRequest) (*response.SendGetUserMeResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "get_user_me",
		MessageData: payload,
		ReplyTo:     "julong_recruitment",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsgPublisher{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendGetUserMe] " + errMsg)
	}

	return &response.SendGetUserMeResponse{
		User: resp.MessageData,
	}, nil
}
