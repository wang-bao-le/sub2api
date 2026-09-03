package admin

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	monitor *service.ConversationMonitorService
}

func NewConversationHandler(monitor *service.ConversationMonitorService) *ConversationHandler {
	return &ConversationHandler{monitor: monitor}
}

func conversationFilter(c *gin.Context) service.ConversationFilter {
	f := service.ConversationFilter{UserSearch: strings.TrimSpace(c.Query("user")), APIKey: strings.TrimSpace(c.Query("api_key")), Model: strings.TrimSpace(c.Query("model")), Provider: strings.TrimSpace(c.Query("provider")), Status: strings.TrimSpace(c.Query("status")), Keyword: strings.TrimSpace(c.Query("keyword")), SortBy: strings.TrimSpace(c.Query("sort_by")), SortOrder: strings.TrimSpace(c.Query("sort_order"))}
	f.Page, f.PageSize = response.ParsePagination(c)
	if raw := c.Query("stream"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			f.Stream = &v
		}
	}
	if raw := c.Query("start_time"); raw != "" {
		if v, err := time.Parse(time.RFC3339, raw); err == nil {
			f.StartTime = &v
		}
	}
	if raw := c.Query("end_time"); raw != "" {
		if v, err := time.Parse(time.RFC3339, raw); err == nil {
			f.EndTime = &v
		}
	}
	return f
}
func (h *ConversationHandler) List(c *gin.Context) {
	page, err := h.monitor.List(c.Request.Context(), conversationFilter(c))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, page.Items, page.Total, page.Page, page.PageSize)
}
func (h *ConversationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid conversation id")
		return
	}
	item, err := h.monitor.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}
func (h *ConversationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid conversation id")
		return
	}
	if err := h.monitor.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
func (h *ConversationHandler) DeleteByFilter(c *gin.Context) {
	n, err := h.monitor.DeleteByFilter(c.Request.Context(), conversationFilter(c))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": n})
}
func (h *ConversationHandler) GetConfig(c *gin.Context) {
	cfg, err := h.monitor.GetConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, cfg)
}
func (h *ConversationHandler) UpdateConfig(c *gin.Context) {
	var cfg service.ConversationMonitorConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid conversation monitor config")
		return
	}
	updated, err := h.monitor.UpdateConfig(c.Request.Context(), cfg)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, updated)
}
