ws业务层协议格式

| 消息指令Command | 消息长度Length | 消息载体Payload |
| --- |------------|-------------|
| 2bytes | 4bytes     | nbytes      |

Command常量定义：
- 100: ping
- 101: pong
- 102: 表示一个文本消息