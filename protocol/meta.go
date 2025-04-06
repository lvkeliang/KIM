package protocol

const (
	// MetaDestServer 消息将要送达的网关的ServiceName
	MetaDestServer = "dest.server"
	// MetaDestChannels 消息将要送达的channels
	MetaDestChannels = "dest.channels"
)

// Command defined data type between client and server
const (
	// login
	CommandLoginSignIn  = "login.signin"
	CommandLoginSignOut = "login.signout"

	// chat
	CommandChatUserTalk  = "chat.user.talk"
	CommandChatGroupTalk = "chat.group.talk"
	CommandChatTalkAck   = "chat.talk.ack"

	// 离线
	CommandOfflineIndex   = "chat.offline.index"
	CommandOfflineContent = "chat.offline.content"

	// 群管理
	CommandGroupCreate  = "chat.group.create"
	CommandGroupJoin    = "chat.group.join"
	CommandGroupQuit    = "chat.group.quit"
	CommandGroupMembers = "chat.group.members"
	CommandGroupDetail  = "chat.group.detail"
)
