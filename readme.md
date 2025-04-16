### 通信层

通信层是对TCP和WebSocket的封装，在inter包中对其进行了接口的定义，实现主要在tcp和websocket包下，将其封装成为了易于收发消息的服务端-客户端形式，可以说是为上层造了个轮子

通信层的协议如下: 

| 消息指令Command | 消息长度Length | 消息载体Payload |
|-------------|------------|-------------|
| 2bytes      | 4bytes     | nbytes      |

Command常量定义:

```go
const (
	OpContinuation OpCode = 0x0
	OpText         OpCode = 0x1
	OpBinary       OpCode = 0x2
	OpClose        OpCode = 0x8
	OpPing         OpCode = 0x9
	OpPong         OpCode = 0xa
)
```

### 基础层

#### 消息协议: 

| 属性         | 类型     | 说明      |
|------------|--------|---------|
| command    | string | 指令      |
| channelId  | string | 连接标识    |
| sequence   | uint32 | 序列号     |
| flag       | enum   | 标识      |
| status     | enum   | 状态码     |
| dest       | string | 目标：群、用户 |
| bodyLength | uint32 | 消息体的长度  |
| body       | uint32 | 消息体  |

command包含如下:

登录类
- login.signin: 登录
- logjin.signout: 退出登录

消息类
- chat.user.talk: 私聊消息
- chat.group.talk: 群聊消息
- chat.talk.ack: ACK

离线类
- chat.offline.index: 下载索引
- chat.offline.content: 下载内容

群操作类
- chat.group.create: 创建群
- chat.group.join: 加入群
- chat.group.quit: 退出群
- chat.group.group.members: 查看群成员

command字段也用于判断dest字段表示的是用户ID还是群聊ID

---

标识flag用于流程控制，有如下值:

- Request: 请求一条消息，通常由客户端发起
- Response：响应一条消息
- Push：一条推送消息

#### 量协议: 

| 消息指令Command | 消息长度Length | 消息载体Payload |
|-------------|------------|-------------|
| 2bytes      | 2bytes     | nbytes      |

轻量协议用于一些类似heartbeat之类的消息

| 协议Command | 说明   |
|-----------|------|
| 1         | ping |
| 2         | pong |

#### 在网关层通过魔数区分两种协议

魔数固定4bytes，封包时可以用反射来进行协议的区分

### 容器层

用于托管Server、维护服务的依赖关系(服务发现和连接建立)、处理消息的上下行

//TODO: 全局时钟(目前是本地时钟)、异步存储(目前是同步)、读写扩散结合(目前是写扩散)