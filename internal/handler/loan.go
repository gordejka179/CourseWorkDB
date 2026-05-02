package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var loanReq struct {
    readerLibraryCard string `json:"readerlibrarycard"`
    copyId int `json:"copyid"`	
}

//делаем выдачу экземпляра
func (h *Handler) makeLoan(c *gin.Context) {
	roleRaw, _ := c.Get("role")

    role, ok := roleRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения role"})
        return
    }

	if role != "librarian"{
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не библиотекарь"})
        return
	}

    emailLibrarianRaw, _ := c.Get("email")

    emailLibrarian, ok := emailLibrarianRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения emailLibrarian"})
        return
    }

    if err := c.BindJSON(&loanReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
        return
    }

    err := h.service.MakeLoan(loanReq.readerLibraryCard, emailLibrarian, loanReq.copyId)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Книга забронирована"})

}

//получения для читателя списка выданных ему экземпляров
func (h *Handler) getLoanedBooks(c *gin.Context) {
	emailReaderRaw, _ := c.Get("email")

    emailReader, ok := emailReaderRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения emailReader"})
        return
    }

	_, _ = h.service.GetLoanedBooks(emailReader)

}
