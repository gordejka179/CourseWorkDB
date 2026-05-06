package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//для выдачи книги логично, что библиотекарь просто отсканирует чит билет и получит его,
//отсканирует qr-код книги и получит её инвентарный номер.
//Но нужно сделать проверку, что у читателя действительно есть такая книга
func (h *Handler) returnBook(c *gin.Context) {
	roleRaw, _ := c.Get("role")
    role, ok := roleRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения role"})
        return
    }

	if role != "librarian"{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы не библиотекарь"})
        return
	}

    switch c.Request.Method {
	case http.MethodGet:
        c.HTML(http.StatusOK, "return.tmpl", gin.H{})
    case http.MethodPost:
        var returnReq struct {
            ReaderLibraryCard string `json:"readerlibrarycard"`
            InventoryNumber string `json:"inventorynumber"`
        }

        if err := c.BindJSON(&returnReq); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
            return
        }

        err := h.service.ReturnBook(returnReq.ReaderLibraryCard, returnReq.InventoryNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Книга возвращена"})

    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не разрешён"})
    }
}