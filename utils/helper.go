package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/IlhamSetiaji/julong-notification-be/internal/request"
	"github.com/IlhamSetiaji/julong-notification-be/internal/response"
	"github.com/gin-gonic/gin"
)

type TemplateHelper struct {
	Ctx *gin.Context
}

func NewTemplateHelper(c *gin.Context) *TemplateHelper {
	return &TemplateHelper{
		Ctx: c,
	}
}

func (h *TemplateHelper) IsAuthenticated() bool {
	return h.Ctx.GetBool("isAuthenticated")
}

func (h *TemplateHelper) NotInArrays(value string, list []string) bool {
	for _, item := range list {
		if value == item {
			return false
		}
	}
	return true
}

func (h *TemplateHelper) CreateSlice(values ...string) []string {
	return values
}

func (h *TemplateHelper) DateFormatter(date time.Time) string {
	return date.Format("2006-01-02")
}

func GenerateRandomIntToken(digits int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%0*d", digits, n.Int64())
	return token, nil
}

func GenerateRandomStringToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

var ResponseChannel = make(chan map[string]interface{}, 100)

var Rchans = make(map[string](chan response.RabbitMQResponse))

type RabbitMsgPublisher struct {
	QueueName string                  `json:"queueName"`
	Message   request.RabbitMQRequest `json:"message"`
}

type RabbitMsgConsumer struct {
	QueueName string                    `json:"queueName"`
	Reply     response.RabbitMQResponse `json:"reply"`
}

// channel to publish rabbit messages
var Pchan = make(chan RabbitMsgPublisher, 10)
var Rchan = make(chan RabbitMsgConsumer, 10)
