package handlers

import (
	"net/http"
	"strconv"
	"social-app/services/chat-service/service"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// 获取聊天房间列表
func (h *ChatHandler) GetChatRooms(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	rooms, err := h.chatService.GetChatRooms(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(rooms))
}

// 创建或获取聊天房间
func (h *ChatHandler) CreateOrGetChatRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.CreateChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	room, err := h.chatService.CreateOrGetChatRoom(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(room))
}

// 获取聊天消息
func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	chatID := c.Param("chatId")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "聊天ID不能为空"))
		return
	}

	// 分页参数
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	messages, err := h.chatService.GetMessages(userID.(string), chatID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(messages))
}

// 发送消息
func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	message, err := h.chatService.SendMessage(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(message))
}

// 标记消息为已读
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	err := h.chatService.MarkAsRead(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "标记成功"}))
}

// 删除消息
func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	messageID := c.Param("messageId")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "消息ID不能为空"))
		return
	}

	err := h.chatService.DeleteMessage(userID.(string), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "删除成功"}))
}

// 上传媒体文件
func (h *ChatHandler) UploadMedia(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	file, err := c.FormFile("media")
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "上传文件错误: "+err.Error()))
		return
	}

	mediaType := c.PostForm("type")
	if mediaType == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "媒体类型不能为空"))
		return
	}

	result, err := h.chatService.UploadMedia(userID.(string), file, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}