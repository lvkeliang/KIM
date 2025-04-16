package rpcService

import (
	"KIM/database"
	"KIM/logger"
	"KIM/naming"
	"KIM/naming/consul"
	"KIM/protocol"
	"KIM/services/rpcService/conf"
	"KIM/services/rpcService/group"
	"KIM/services/rpcService/offline"
	"KIM/services/rpcService/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type ServerStartOptions struct {
	Config  string
	RunType string
}

func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	initConf, err := conf.Init(opts.Config)
	if err != nil {
		return err
	}
	_ = logger.Init(logger.Settings{
		Level: "info",
	})

	// database.Init
	var (
		baseDb    *gorm.DB
		messageDb *gorm.DB
	)
	baseDb, err = database.InitMysqlDb(initConf.Driver, initConf.BaseDb)
	if err != nil {
		return err
	}
	messageDb, err = database.InitMysqlDb(initConf.Driver, initConf.MessageDb)
	if err != nil {
		return err
	}

	_ = baseDb.AutoMigrate(&database.Group{}, &database.GroupMember{})
	_ = messageDb.AutoMigrate(&database.MessageIndex{}, &database.MessageContent{})

	if initConf.NodeID == 0 {
		initConf.NodeID = int64(util.HashCode(initConf.ServiceID))
	}
	idgen, err := database.NewIDGenerator(initConf.NodeID)
	if err != nil {
		return err
	}

	rdb, err := conf.InitRedis(initConf.RedisAddrs, "")
	if err != nil {
		return err
	}

	ns, err := consul.NewNaming(initConf.ConsulURL)
	if err != nil {
		return err
	}

	service := &naming.DefaultServiceRegistration{
		Id:       initConf.ServiceID,
		Name:     protocol.SNService, // service name
		Address:  initConf.PublicAddress,
		Port:     initConf.PublicPort,
		Protocol: "http",
		Tags:     initConf.Tags,
		Meta: map[string]string{
			consul.KeyHealthURL: fmt.Sprintf("http://%s:%d/health", initConf.PublicAddress, initConf.PublicPort),
		},
	}

	if opts.RunType == "local" {
		service.Address = "127.0.0.1"
		service.Meta[consul.KeyHealthURL] = fmt.Sprintf("http://%s:%d/health", "127.0.0.1", initConf.PublicPort)
	} else if opts.RunType == "docker" {
		service.Address = "127.0.0.1"
		service.Meta[consul.KeyHealthURL] = fmt.Sprintf("http://%s:%d/health", "host.docker.internal", initConf.PublicPort)
	}

	_ = ns.Register(service)
	defer func() {
		_ = ns.Deregister(initConf.ServiceID)
	}()
	serviceHandlerOffline := offline.ServiceHandlerOffline{
		BaseDb:    baseDb,
		MessageDb: messageDb,
		Idgen:     idgen,
		Cache:     rdb,
	}

	serviceHandlerGroup := group.ServiceHandlerGroup{
		BaseDb:    baseDb,
		MessageDb: messageDb,
		Idgen:     idgen,
		Cache:     rdb,
	}

	ac := util.MakeAccessLog()
	defer ac.Close()

	app := newApp(&serviceHandlerOffline, &serviceHandlerGroup)

	app.Use(ac.Middleware())
	app.Use(negotiationMiddleware())

	return startServerWithOptimizations(app, ":8080")
}

func newApp(serviceHandlerOffline *offline.ServiceHandlerOffline, serviceHandlerGroup *group.ServiceHandlerGroup) *gin.Engine {
	app := gin.Default()

	app.Use(CORSMiddleware())

	app.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	messageAPI := app.Group("/api/:app/message")
	{

		messageAPI.POST("/user", serviceHandlerOffline.InsertUserMessage)   // 离线服务，向离线队列添加一条单聊消息
		messageAPI.POST("/group", serviceHandlerOffline.InsertGroupMessage) // 离线服务，向离线队列添加一条群聊消息
		messageAPI.POST("/ack", serviceHandlerOffline.MessageAck)           // 离线服务，通过ACK确认接收到消息，只用确认接受到的最后一条消息，表示该消息及其之前的消息都已经确认
	}

	groupAPI := app.Group("/api/:app/group")
	{
		groupAPI.GET("/:id", serviceHandlerGroup.GroupGet)
		groupAPI.POST("", serviceHandlerGroup.GroupCreate)
		groupAPI.POST("/member", serviceHandlerGroup.GroupJoin)
		groupAPI.DELETE("/member", serviceHandlerGroup.GroupQuit)
		groupAPI.GET("/members/:id", serviceHandlerGroup.GroupMembers)

	}

	offlineAPI := app.Group("/api/:app/offline")
	{
		offlineAPI.Use(util.AutoCompression())
		offlineAPI.POST("/index", serviceHandlerOffline.GetOfflineMessageIndex)     //获取消息索引
		offlineAPI.POST("/content", serviceHandlerOffline.GetOfflineMessageContent) //下载消息内容
	}
	return app
}

// CORSMiddleware 自定义CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求（OPTIONS）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func negotiationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.GetHeader("Accept")

		// 支持的格式
		supported := []string{
			"application/json",
			"application/x-protobuf",
			"application/x-msgpack",
		}

		// 默认JSON
		if accept == "" || accept == "*/*" {
			c.Request.Header.Set("Accept", "application/json")
			c.Next()
			return
		}

		// 检查是否支持请求的格式
		for _, mime := range supported {
			if strings.Contains(accept, mime) {
				c.Next()
				return
			}
		}

		// 不支持的格式返回406
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error":     "Unsupported response format",
			"supported": supported,
		})
	}
}

func startServerWithOptimizations(router *gin.Engine, listenAddr string) error {
	// 1. 设置发布模式（禁用调试信息）
	gin.SetMode(gin.ReleaseMode)

	// 2. 禁用控制台颜色（减少输出处理）
	gin.DisableConsoleColor()

	// 3. 创建自定义 HTTP 服务器
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: router,

		// 性能优化相关配置
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB

		// 更多优化选项
		ErrorLog: log.New(io.Discard, "", 0), // 禁用错误日志
	}

	//// 4. 启动前预热（可选）
	//go func() {
	//	time.Sleep(100 * time.Millisecond)
	//	http.Get("http://" + listenAddr + "/health")
	//}()

	// 5. 启动服务器
	fmt.Printf("Server listening on %s\n", listenAddr)
	return srv.ListenAndServe()
}
