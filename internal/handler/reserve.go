package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var reserveReq struct {
    CopyId string `json:"copyid"`
}

//сделать бронирование
func (h *Handler) reserve(c *gin.Context) {
    emailRaw, _ := c.Get("email")

    email, ok := emailRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения email"})
        return
    }

    if err := c.BindJSON(&reserveReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
        return
    }

	copyID, _ := strconv.Atoi(reserveReq.CopyId)

    err := h.service.ReserveCopyByEmail(email, copyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Книга забронирована"})

}

