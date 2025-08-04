package types

// WebSocket消息类型
type WSMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	UserID  string      `json:"user_id,omitempty"`
	ChatID  string      `json:"chat_id,omitempty"`
}

// WebSocket消息类型常量
const (
	WSMessageTypeNewMessage    = "new_message"
	WSMessageTypeMessageRead   = "message_read"
	WSMessageTypeTyping        = "typing"
	WSMessageTypeStopTyping    = "stop_typing"
	WSMessageTypeUserOnline    = "user_online"
	WSMessageTypeUserOffline   = "user_offline"
	WSMessageTypeCallRequest   = "call_request"
	WSMessageTypeCallResponse  = "call_response"
	WSMessageTypeCallEnd       = "call_end"
)