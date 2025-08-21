package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"social-app/services/chat-service/types"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域，生产环境应该限制
	},
}

// Client 表示一个WebSocket客户端
type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

// Hub 管理所有WebSocket连接
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool
	
	// 用户ID到客户端的映射
	userClients map[string]*Client
	
	// 注册请求
	register chan *Client
	
	// 注销请求
	unregister chan *Client
	
	// 广播消息
	broadcast chan []byte
	
	// 发送给特定用户
	sendToUser chan UserMessage
	
	// 互斥锁
	mutex sync.RWMutex
}

// UserMessage 发送给特定用户的消息
type UserMessage struct {
	UserID  string
	Message []byte
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[string]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		sendToUser:  make(chan UserMessage),
	}
}

// Run 启动Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.userClients[client.UserID] = client
			h.mutex.Unlock()
			
			log.Printf("User %s connected", client.UserID)
			
			// 通知其他用户该用户上线
			onlineMsg := types.WSMessage{
				Type:   types.WSMessageTypeUserOnline,
				UserID: client.UserID,
			}
			if msgBytes, err := json.Marshal(onlineMsg); err == nil {
				h.broadcast <- msgBytes
			}

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userClients, client.UserID)
				close(client.Send)
			}
			h.mutex.Unlock()
			
			log.Printf("User %s disconnected", client.UserID)
			
			// 通知其他用户该用户下线
			offlineMsg := types.WSMessage{
				Type:   types.WSMessageTypeUserOffline,
				UserID: client.UserID,
			}
			if msgBytes, err := json.Marshal(offlineMsg); err == nil {
				h.broadcast <- msgBytes
			}

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
					delete(h.userClients, client.UserID)
				}
			}
			h.mutex.RUnlock()

		case userMsg := <-h.sendToUser:
			h.mutex.RLock()
			if client, ok := h.userClients[userMsg.UserID]; ok {
				select {
				case client.Send <- userMsg.Message:
				default:
					close(client.Send)
					delete(h.clients, client)
					delete(h.userClients, client.UserID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// SendToUser 发送消息给特定用户
func (h *Hub) SendToUser(userID string, message []byte) {
	h.sendToUser <- UserMessage{
		UserID:  userID,
		Message: message,
	}
}

// Broadcast 广播消息给所有用户
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// GetOnlineUsers 获取在线用户列表
func (h *Hub) GetOnlineUsers() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	users := make([]string, 0, len(h.userClients))
	for userID := range h.userClients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID string) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	_, exists := h.userClients[userID]
	return exists
}

// HandleWebSocket 处理WebSocket连接
func (h *Hub) HandleWebSocket(c *gin.Context) {
	// 从查询参数或JWT token中获取用户ID
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:     generateClientID(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h,
	}

	// 注册客户端
	h.register <- client

	// 启动goroutines
	go client.writePump()
	go client.readPump()
}

// readPump 处理从WebSocket连接读取消息
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 处理接收到的消息
		var wsMsg types.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// 设置发送者ID
		wsMsg.UserID = c.UserID

		// 根据消息类型处理
		switch wsMsg.Type {
		case types.WSMessageTypeTyping:
			// 转发打字状态给聊天室其他成员
			if wsMsg.ChatID != "" {
				c.forwardToChatMembers(wsMsg, message)
			}
		case types.WSMessageTypeStopTyping:
			// 转发停止打字状态
			if wsMsg.ChatID != "" {
				c.forwardToChatMembers(wsMsg, message)
			}
		default:
			log.Printf("Unknown message type: %s", wsMsg.Type)
		}
	}
}

// writePump 处理向WebSocket连接写入消息
func (c *Client) writePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// forwardToChatMembers 转发消息给聊天室其他成员
func (c *Client) forwardToChatMembers(wsMsg types.WSMessage, message []byte) {
	// 获取聊天室参与者列表
	c.Hub.mutex.RLock()
	defer c.Hub.mutex.RUnlock()
	
	// 转发给聊天室的其他成员
	for userID, client := range c.Hub.userClients {
		if userID != c.UserID { // 不发送给自己
			select {
			case client.Send <- message:
			default:
				// 如果发送失败，关闭连接
				close(client.Send)
				delete(c.Hub.clients, client)
				delete(c.Hub.userClients, userID)
			}
		}
	}
}

// generateClientID 生成客户端ID
func generateClientID() string {
	// 简单的ID生成，实际应该使用UUID
	return "client_" + string(rune(len("temp")))
}