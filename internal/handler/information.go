package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *Handler) getCurrentBookings(c *gin.Context) {
    emailRaw, _ := c.Get("email")

    email, ok := emailRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения email"})
        return
    }

    bi, err := h.service.GetCurrentBookingsByEmail(email)

    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }


    c.HTML(http.StatusOK, "bookingsReader.tmpl", bi)

}