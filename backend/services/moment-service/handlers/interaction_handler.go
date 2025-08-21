package handlers

import (
	"net/http"
	"social-app/services/moment-service/repository"
	"social-app/services/moment-service/service"
	"social-app/shared/models"

	"github.com/gin-gonic/gin"
)

type InteractionHandler struct {
	interactionRepo repository.InteractionRepository
}

func NewInteractionHandler(interactionRepo repository.InteractionRepository) *InteractionHandler {
	return &InteractionHandler{
		interactionRepo: interactionRepo,
	}
}

// 点赞/取消点赞
func (h *InteractionHandler) ToggleLike(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	momentID := c.Param("moment_id")
	if momentID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "动态ID不能为空"))
		return
	}

	// 检查是否已点赞
	isLiked, err := h.interactionRepo.IsLiked(momentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "操作失败"))
		return
	}

	if isLiked {
		// 取消点赞
		err = h.interactionRepo.UnlikeMoment(momentID, userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "取消点赞失败"))
			return
		}
		c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"liked": false, "message": "取消点赞成功"}))
	} else {
		// 点赞
		err = h.interactionRepo.LikeMoment(momentID, userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "点赞失败"))
			return
		}
		c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"liked": true, "message": "点赞成功"}))
	}
}

// 获取动态点赞列表
func (h *InteractionHandler) GetLikes(c *gin.Context) {
	momentID := c.Param("moment_id")
	if momentID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "动态ID不能为空"))
		return
	}

	likes, err := h.interactionRepo.GetMomentLikes(momentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "获取点赞列表失败"))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(likes))
}

// 添加评论
func (h *InteractionHandler) AddComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	momentID := c.Param("moment_id")
	if momentID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "动态ID不能为空"))
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "评论内容不能为空"))
		return
	}

	comment := &models.MomentComment{
		MomentID: momentID,
		UserID:   userID.(string),
		Content:  req.Content,
	}

	err := h.interactionRepo.CreateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "评论失败"))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "评论成功"}))
}

// 获取动态评论列表
func (h *InteractionHandler) GetComments(c *gin.Context) {
	momentID := c.Param("moment_id")
	if momentID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "动态ID不能为空"))
		return
	}

	comments, err := h.interactionRepo.GetMomentComments(momentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "获取评论列表失败"))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(comments))
}