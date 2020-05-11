package middleware

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/sys520084/namenode/internal/log"
)

var timeFormat = "2006-01-02 15:04:05.99"

// gin日志中间件
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {

		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewV1().String()
		}

		log.SetApiRequestID(c, requestID)

		// 获取api的请求id
		// log.GetApiRequestID(c)

		baseEntry := log.UniqApiLogEntry(requestID)
		path := c.Request.URL.Path

		var body []byte
		if c.Request.Body != nil {
			body, err := ioutil.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
			}
		}

		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		baseEntry = baseEntry.WithFields(logrus.Fields{
			"time":      time.Now().Format(timeFormat),
			"path":      path,
			"method":    c.Request.Method,
			"host":      c.Request.Host,
			"raw_query": c.Request.URL.RawQuery,
			"userAgent": clientUserAgent,
			"referer":   referer,
			"hostname":  hostname,
			"clientIP":  clientIP,
		})
		if len(body) != 0 {
			baseEntry = baseEntry.WithField("body", string(body))
		}
		baseEntry.Info("http")
		defer log.DeferLog(baseEntry)

		c.Next()

		// api结束时的操作
	}

}
