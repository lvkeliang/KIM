package util

import (
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"hash/crc32"
	"net/http"
	"strings"
)

// Negotiate 自动编码并发送
func Negotiate(c *gin.Context, status int, data interface{}) {
	accept := c.GetHeader("Accept")

	switch {
	case strings.Contains(accept, "application/json"):
		c.JSON(status, data)
	case strings.Contains(accept, "text/xml"):
		c.XML(status, data)
	case strings.Contains(accept, "application/protobuf"),
		strings.Contains(accept, "application/x-protobuf"):

		// 类型断言：检查 data 是否实现了 proto.Message
		msg, ok := data.(proto.Message)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to marshal protobuf",
			})
		}
		c.ProtoBuf(status, msg)
	default:
		c.JSON(status, data) // 默认 JSON
	}
}

// AutoBind 自动解析Body内容
func AutoBind(c *gin.Context, data interface{}) error {
	contentType := c.ContentType()

	switch {
	case strings.Contains(contentType, "application/json"):
		return c.ShouldBindJSON(data)
	case strings.Contains(contentType, "application/xml"):
		return c.ShouldBindXML(data)
	case strings.Contains(contentType, "application/protobuf"),
		strings.Contains(contentType, "application/x-protobuf"):
		// 类型断言：检查 data 是否实现了 proto.Message
		if msg, ok := data.(proto.Message); ok {
			return bindProtobuf(c, msg) // 安全调用
		}
		// 类型不匹配时返回明确错误
		return fmt.Errorf("data must implement proto got type %T", data)
	default:
		return c.ShouldBind(data) // 尝试通用绑定
	}
}

// ErrorResponse 统一错误处理
func ErrorResponse(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
}

//func protobufResponse(c *gin.Context, status int, data proto.Message) {
//	bytes, err := proto.Marshal(data)
//	if err != nil {
//		c.AbortWithStatus(http.StatusInternalServerError)
//		return
//	}
//	c.Data(status, "application/x-protobuf", bytes)
//}

func bindProtobuf(c *gin.Context, data proto.Message) error {
	bytes, err := c.GetRawData()
	if err != nil {
		return err
	}
	return proto.Unmarshal(bytes, data)
}

// CompressionType 压缩类型
type CompressionType int

const (
	Gzip CompressionType = iota
	Brotli
	None
)

// AutoCompression 自动使用压缩算法中间件
func AutoCompression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查客户端支持的压缩类型
		acceptEncoding := c.GetHeader("Accept-Encoding")
		compressionType := getPreferredCompression(acceptEncoding)

		switch compressionType {
		case Gzip:
			gzip.Gzip(gzip.DefaultCompression)(c)
		case Brotli:
			brotliMiddleware()(c)
		case None:
			c.Next()
		}
	}
}

// 根据 Accept-Encoding 头选择最佳压缩算法
func getPreferredCompression(acceptEncoding string) CompressionType {
	encodings := strings.Split(acceptEncoding, ",")

	for _, enc := range encodings {
		enc = strings.TrimSpace(enc)
		if strings.Contains(enc, "br") {
			return Brotli
		}
		if strings.Contains(enc, "gzip") {
			return Gzip
		}
	}

	return None
}

// Brotli 压缩中间件
func brotliMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过已经压缩的内容或小文件
		if shouldSkipCompression(c) {
			c.Next()
			return
		}

		// 设置响应头
		c.Header("Content-Encoding", "br")
		c.Header("Vary", "Accept-Encoding")

		// 使用 Brotli 压缩
		brWriter := brotli.NewWriter(c.Writer)
		defer brWriter.Close()

		c.Writer = &compressedResponseWriter{
			ResponseWriter: c.Writer,
			writer:         brWriter,
		}

		c.Next()
	}
}

// 判断是否应该跳过压缩
func shouldSkipCompression(c *gin.Context) bool {
	// 跳过已经压缩的内容
	if c.GetHeader("Content-Encoding") != "" {
		return true
	}

	// 跳过特定内容类型
	contentType := c.Writer.Header().Get("Content-Type")
	switch {
	case strings.HasPrefix(contentType, "image/"),
		strings.HasPrefix(contentType, "video/"),
		strings.HasPrefix(contentType, "audio/"),
		contentType == "application/pdf":
		return true
	}

	return false
}

// 压缩响应写入器
type compressedResponseWriter struct {
	gin.ResponseWriter
	writer *brotli.Writer
}

func (w *compressedResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func HashCode(key string) uint32 {
	hash32 := crc32.NewIEEE()
	hash32.Write([]byte(key))
	return hash32.Sum32() % 1000
}
