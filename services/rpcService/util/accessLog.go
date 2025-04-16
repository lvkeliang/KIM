package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLogConfig struct {
	Delim              string
	TimeFormat         string
	Async              bool
	IP                 bool
	BytesReceived      bool
	BytesSent          bool
	BytesReceivedBody  bool
	BytesSentBody      bool
	BodyMinify         bool
	RequestBody        bool
	ResponseBody       bool
	KeepMultiLineError bool
	Output             io.Writer
}

type AccessLog struct {
	config    *AccessLogConfig
	file      *os.File
	logChan   chan string
	closeChan chan struct{}
	wg        sync.WaitGroup
	mu        sync.Mutex
}

func NewAccessLog(config *AccessLogConfig) *AccessLog {
	// 默认配置
	if config.Delim == "" {
		config.Delim = "|"
	}
	if config.TimeFormat == "" {
		config.TimeFormat = "2006-01-02 15:04:05"
	}

	// 创建日志文件
	file, err := os.OpenFile("./access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	al := &AccessLog{
		config:    config,
		file:      file,
		logChan:   make(chan string, 1000), // 缓冲通道
		closeChan: make(chan struct{}),
	}

	if al.config.Async {
		al.wg.Add(1)
		go al.asyncWriter()
	}

	return al
}

func (a *AccessLog) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 包装ResponseWriter以跟踪响应大小
		writer := &responseWriterWrapper{
			ResponseWriter: c.Writer,
			config:         a.config,
		}
		c.Writer = writer

		// 记录请求开始时间
		start := time.Now()

		// 读取请求体
		var reqBody []byte
		if a.config.RequestBody && c.Request.Body != nil {
			reqBody, _ = c.GetRawData()
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 构建日志条目
		logEntry := a.buildLogEntry(c, writer, reqBody, latency)

		// 写入日志
		if a.config.Async {
			select {
			case a.logChan <- logEntry:
			default:
				// 如果通道满了，丢弃日志条目避免阻塞
				fmt.Println("Log channel full, dropping log entry")
			}
		} else {
			a.writeLog(logEntry)
		}
	}
}

func (a *AccessLog) asyncWriter() {
	defer a.wg.Done()

	for {
		select {
		case entry := <-a.logChan:
			a.writeLog(entry)
		case <-a.closeChan:
			// 处理剩余日志条目
			for {
				select {
				case entry := <-a.logChan:
					a.writeLog(entry)
				default:
					return
				}
			}
		}
	}
}

func (a *AccessLog) Close() {
	if a.config.Async {
		close(a.closeChan)
		a.wg.Wait()
	}
	a.file.Close()
}

func (a *AccessLog) buildLogEntry(c *gin.Context, writer *responseWriterWrapper, reqBody []byte, latency time.Duration) string {
	var builder strings.Builder

	// 时间戳
	builder.WriteString(time.Now().Format(a.config.TimeFormat))
	builder.WriteString(a.config.Delim)

	// 延迟
	builder.WriteString(latency.String())
	builder.WriteString(a.config.Delim)

	// 状态码
	builder.WriteString(fmt.Sprintf("%d", writer.Status()))
	builder.WriteString(a.config.Delim)

	// 方法
	builder.WriteString(c.Request.Method)
	builder.WriteString(a.config.Delim)

	// 路径
	builder.WriteString(c.Request.URL.Path)
	builder.WriteString(a.config.Delim)

	// IP地址
	if a.config.IP {
		builder.WriteString(c.ClientIP())
	}
	builder.WriteString(a.config.Delim)

	// 查询参数
	builder.WriteString(c.Request.URL.RawQuery)
	builder.WriteString(a.config.Delim)

	// 接收字节数
	if a.config.BytesReceived {
		builder.WriteString(fmt.Sprintf("%d", c.Request.ContentLength))
	}
	builder.WriteString(a.config.Delim)

	// 发送字节数
	if a.config.BytesSent {
		builder.WriteString(fmt.Sprintf("%d", writer.size))
	}
	builder.WriteString(a.config.Delim)

	// 请求体字节数
	if a.config.BytesReceivedBody {
		builder.WriteString(fmt.Sprintf("%d", len(reqBody)))
	}
	builder.WriteString(a.config.Delim)

	// 响应体字节数
	if a.config.BytesSentBody {
		builder.WriteString(fmt.Sprintf("%d", writer.size))
	}
	builder.WriteString(a.config.Delim)

	// 请求体内容
	if a.config.RequestBody {
		body := string(reqBody)
		if a.config.BodyMinify {
			body = minifyBody(body)
		}
		builder.WriteString(body)
	}
	builder.WriteString(a.config.Delim)

	// 响应体内容
	if a.config.ResponseBody {
		builder.WriteString("(response body not captured)")
	}
	builder.WriteString(a.config.Delim)

	// 错误信息
	errors := c.Errors.String()
	if a.config.KeepMultiLineError {
		errors = strings.ReplaceAll(errors, "\n", "\\n")
	} else {
		if idx := strings.Index(errors, "\n"); idx != -1 {
			errors = errors[:idx] + "..."
		}
	}
	builder.WriteString(errors)

	builder.WriteString("\n")

	return builder.String()
}

func (a *AccessLog) writeLog(entry string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.config.Output != nil {
		a.config.Output.Write([]byte(entry))
	}
	a.file.WriteString(entry)
}

func minifyBody(body string) string {
	body = strings.ReplaceAll(body, "\n", "")
	body = strings.ReplaceAll(body, "\t", "")
	body = strings.ReplaceAll(body, " ", "")
	return body
}

type responseWriterWrapper struct {
	gin.ResponseWriter
	size   int
	config *AccessLogConfig
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *responseWriterWrapper) WriteString(s string) (int, error) {
	size, err := w.ResponseWriter.WriteString(s)
	w.size += size
	return size, err
}

func MakeAccessLog() *AccessLog {
	// 配置访问日志
	config := &AccessLogConfig{
		Delim:              "|",
		TimeFormat:         "2006-01-02 15:04:05",
		Async:              true, // 默认启用异步
		IP:                 true,
		BytesReceived:      false,
		BytesSent:          false,
		BytesReceivedBody:  true,
		BytesSentBody:      true,
		BodyMinify:         true,
		RequestBody:        true,
		ResponseBody:       false,
		KeepMultiLineError: true,
		Output:             os.Stdout,
	}

	return NewAccessLog(config)
}
