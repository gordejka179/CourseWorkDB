package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) search(c *gin.Context) {
	c.HTML(http.StatusOK, "search.tmpl", gin.H{})
}