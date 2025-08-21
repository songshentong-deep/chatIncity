package handlers

import (
	"net/http"
	"social-app/services/moment-service/service"

	"github.com/gin-gonic/gin"
)

type MomentHandler struct {
	momentService service.MomentService
}

func NewMomentHandler(momentService service.MomentService) *MomentHandler {
	return &MomentHandler{
		momentService: momentService,
	}
}

// 创建动态
func (h *MomentHandler) CreateMoment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.CreateMomentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	req.UserID = userID.(string)

	err := h.momentService.CreateMoment(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "动态发布成功"}))
}

// 获取我的动态
func (h *MomentHandler) GetMyMoments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	moments, err := h.momentService.GetMyMoments(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(moments))
}

// 获取好友动态
func (h *MomentHandler) GetFriendsMoments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	moments, err := h.momentService.GetFriendsMoments(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(moments))
}

// 删除动态
func (h *MomentHandler) DeleteMoment(c *gin.Context) {
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

	err := h.momentService.DeleteMoment(momentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "动态删除成功"}))
}