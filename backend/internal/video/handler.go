package video

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/idcim/aihuman/backend/internal/api"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "invalid request body")
		return
	}

	task, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	api.Success(c, task)
}

func (h *Handler) Get(c *gin.Context) {
	task, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if errors.Is(err, ErrTaskNotFound) {
		api.NotFound(c, "video task not found")
		return
	}
	if err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	api.Success(c, task)
}
