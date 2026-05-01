package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var req struct {
    CopyId string `json:"copyid"`
}

func (h *Handler) reserve(c *gin.Context) {
    emailRaw, _ := c.Get("email")

    email, ok := emailRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения email"})
        return
    }

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
        return
    }

	copyID, _ := strconv.Atoi(req.CopyId)

    err := h.service.ReserveCopyByEmail(email, copyID)

    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Книга забронирована"})

}