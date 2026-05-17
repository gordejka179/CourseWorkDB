package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OverallStats struct {
    Total int `json:"total"`
    Available int `json:"available"`
    Reserved int `json:"reserved"`
    LoanedOut int `json:"loaned_out"`
    Overdue int `json:"overdue"`
}

func (h *Handler) getInfo(c *gin.Context) {
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
        c.HTML(http.StatusOK, "searchLibrarianInfo.tmpl", gin.H{})
	}
}

func (h *Handler) overdue(c *gin.Context){
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
        overdueCopies, err := h.service.GetOverdueCopies()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка вызова GetOverdueCopies"})
		}
		c.JSON(http.StatusOK, overdueCopies)
		return
	}
}


func (h *Handler) overall(c *gin.Context){
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
    	stats, err := h.service.GetOverallStats()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка вызова GetOverdueCopies"})
		}
		c.JSON(http.StatusOK, stats)
		return
	}
}
