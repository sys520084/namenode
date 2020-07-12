package status

import (
	"github.com/gin-gonic/gin"
	"github.com/sys520084/namenode/internal/log"
)

const (
	SuccessCode      int = 0
	NotLoginCode     int = 10000
	NoPrivCode       int = 10001
	BadRequestCode   int = 10002
	HandlerErrorCode int = 10003
)

type Status struct {
	Code      int         `json:"code"`
	RequestID string      `json:"request_id"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
}

// 200
func StatusOK(c *gin.Context, msg string, Data interface{}) Status {
	requestID := log.GetApiRequestID(c)
	return _build(SuccessCode, msg, requestID, Data)
}

// 10000
func NotLoginStatus(c *gin.Context, msg string, Data interface{}) Status {
	requestID := log.GetApiRequestID(c)
	return _build(NotLoginCode, msg, requestID, Data)
}

// 10001
func NoPrivStatus(c *gin.Context, msg string, Data interface{}) Status {
	requestID := log.GetApiRequestID(c)
	return _build(NoPrivCode, msg, requestID, Data)
}

// 10002
func BadRequestStatus(c *gin.Context, msg string, Data interface{}) Status {
	requestID := log.GetApiRequestID(c)
	return _build(BadRequestCode, msg, requestID, Data)
}

// 10003
func HandlerErrorStatus(c *gin.Context, msg string, Data interface{}) Status {
	requestID := log.GetApiRequestID(c)
	return _build(HandlerErrorCode, msg, requestID, Data)
}

func _build(code int, msg, requestID string, Data interface{}) Status {
	return Status{
		RequestID: requestID,
		Code:      code,
		Msg:       msg,
		Data:      Data,
	}
}

func MaybePanic(err error) {
	if err != nil {
		panic(Status{Msg: err.Error()})
	}
}

func Panic(msg string) {
	panic(Status{Msg: msg})
}
