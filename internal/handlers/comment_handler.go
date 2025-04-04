package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmfadel/wander-base/internal/models"
	"github.com/wmfadel/wander-base/internal/models/core"
	"github.com/wmfadel/wander-base/internal/service"
	"github.com/wmfadel/wander-base/pkg/utils"
)

type CommentHandler struct {
	CommentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{CommentService: commentService}
}

func (h *CommentHandler) Create(c *gin.Context) {
	var commentRequest models.CreateCommentRequest
	if err := c.ShouldBindJSON(&commentRequest); err != nil {
		c.JSON(http.StatusBadRequest, core.NewESError("Failed to parse comment", err))
		return
	}

	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewESError("Failed to parse event ID", err))
		return
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, core.NewESError("Failed to get user from context", err))
		return
	}

	comment := models.Comment{
		EventID: eventId,
		UserID:  user.ID,
		Content: commentRequest.Content,
		Score:   0,
		Visible: true,
	}
	err = h.CommentService.Create(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.NewESError("Failed to create comment", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created"})
}

func (h *CommentHandler) GetEventComments(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewESError("Failed to parse event ID", err))
		return
	}
	comments, err := h.CommentService.GetEventComments(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.NewESError("Failed to get comments", err))
		return
	}
	c.JSON(http.StatusOK, comments)
}
