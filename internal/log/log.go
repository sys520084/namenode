package log

import (
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LOG_KEY        string = "LOG_KEY"
	REQUEST_ID_KEY string = "REQUEST_ID_KEY"
)

var Logger *logrus.Logger

//  初始化日志对象
func InitLogger() {

	Logger = logrus.New()

	// 设置输出的日志等级
	if viper.GetBool("debug") {
		Logger.SetLevel(logrus.DebugLevel)
	}

	// 设置输出的日志等级
	if viper.GetBool("log_json_format") {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	logfile := strings.TrimSpace(viper.GetString("logfile"))

	// 有文件路径则写入文件
	if logfile != "" {
		dir := filepath.Dir(logfile)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			Logger.Errorf("'%s' logfile dir '%s' IsNotExist", logfile, dir)
			return
		}
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			Logger.Errorf("'%s' logfile os.OpenFile failed:%s", logfile, dir, err)
			return
		}
		Logger.SetOutput(file)
	}

}

// 提供api唯一ID
func UniqApiLogEntry(requestID string) *logrus.Entry {
	entry := logrus.NewEntry(Logger).WithFields(logrus.Fields{
		"request_id": requestID,
		"action":     "httpApi",
	})
	return entry
}

// 获取日志对象
func GetApiLogEntry(c *gin.Context) *logrus.Entry {
	v, exists := c.Get(LOG_KEY)
	if !exists {
		requestID := uuid.NewV1().String()
		return UniqApiLogEntry(requestID)
	}
	entry, ok := v.(*logrus.Entry)
	if !ok {
		requestID := uuid.NewV1().String()
		return UniqApiLogEntry(requestID)

	}
	return entry
}

// 获取唯一ID
func GetApiRequestID(c *gin.Context) string {
	v, exists := c.Get(REQUEST_ID_KEY)
	if !exists {
		return ""
	}
	requestID, _ := v.(string)
	return requestID
}

// 给gin.Context写入请求ID
func SetApiRequestID(c *gin.Context, v string) {
	c.Set(REQUEST_ID_KEY, v)
}

// 结束及异常处理
func DeferLog(entry *logrus.Entry) {
	nowTime := time.Now()
	entry = entry.WithFields(logrus.Fields{
		"end_time": nowTime,
	})
	if err := recover(); err != nil {
		fileName, line := fileStackLog()
		entry = entry.WithFields(logrus.Fields{
			"exec_line": line,
			"file_name": fileName,
		})
		entry.WithTime(nowTime).Error(err)
		return
	}
	entry.WithTime(nowTime).Debugln("end")
}

// 栈日志
func fileStackLog() (string, int) {
	var name string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		_, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}
	if name != "" {
		return name, line
	}

	return "???", 0
}

func latency(start time.Time) int {
	return int(math.Ceil(float64(time.Since(start).Nanoseconds()) / 1000000.0))
}
