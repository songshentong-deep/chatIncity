package handlers

import (
	"net/http"
	"social-app/services/news-service/models"
	"social-app/services/news-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	newsService service.NewsService
}

func NewNewsHandler(newsService service.NewsService) *NewsHandler {
	return &NewsHandler{
		newsService: newsService,
	}
}

// 获取新闻列表
func (h *NewsHandler) GetNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	news, err := h.newsService.GetNews(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取新闻失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"news":  news,
			"page":  page,
			"limit": limit,
		},
	})
}

// 获取分类列表
func (h *NewsHandler) GetCategories(c *gin.Context) {
	categories, err := h.newsService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取分类失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    categories,
	})
}

// 根据分类获取新闻
func (h *NewsHandler) GetNewsByCategory(c *gin.Context) {
	categoryStr := c.Param("category")
	category := models.NewsCategory(categoryStr)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	news, err := h.newsService.GetNewsByCategory(category, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取分类新闻失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"news":     news,
			"category": category,
			"page":     page,
			"limit":    limit,
		},
	})
}

// 获取新闻详情
func (h *NewsHandler) GetNewsDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "新闻ID不能为空",
		})
		return
	}

	news, err := h.newsService.GetNewsDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "新闻不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    news,
	})
}

// 刷新新闻（管理员接口）
func (h *NewsHandler) RefreshNews(c *gin.Context) {
	err := h.newsService.RefreshNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "刷新新闻失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "新闻刷新成功",
	})
}