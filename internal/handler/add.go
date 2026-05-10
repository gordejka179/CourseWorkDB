package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type checkForm struct {
	Query string `json:"query"`
}

type createAuthorForm struct {
	LastName string `json:"lastName"`
    FirstName string `json:"firstName"`
	Patronymic string `json:"patronymic"`
	BirthDate string `json:"birthDate"`
}

type createPublicationForm struct {
    ISBN string `json:"isbn"`
    OtherIsbn string `json:"otherIsbn"`
    OtherIndexes string `json:"otherIndexes"`
    BBK string `json:"bbk"`
    Title string `json:"title"`
    Year string `json:"year"`
    AuthorIds []int `json:"authorIds"`
}

type authorResponse struct {
    AuthorId int `json:"id"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Patronymic string `json:"patronymic"`
    BirthDate string `json:"birthDate"`
}


func (h *Handler) checkAuthor(c *gin.Context) {
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
    case http.MethodPost:
        var form checkForm
        if err := c.BindJSON(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
            return
        }
		parts := strings.Split(form.Query, "|")
		authors, err := h.service.SearchAuthors(parts[0], parts[1], parts[2], parts[3])
        if err != nil{
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске авторов"})
        }
        
        response := make([]authorResponse, len(authors))

        for i, a := range authors {
            response[i] = authorResponse{
                AuthorId:   a.AuthorId,
                FirstName:  a.FirstName,
                LastName:   a.LastName,
                Patronymic: a.Patronymic,
                BirthDate:  a.BirthDate,
            }
        }

        c.JSON(http.StatusOK, response)

    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не разрешён"})
    }
}

//Создать автора
func (h *Handler) createAuthor(c *gin.Context) {
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
    case http.MethodPost:
        var form createAuthorForm
        if err := c.BindJSON(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
            return
        }
		err := h.service.CreateAuthor(form.LastName, form.FirstName, form.Patronymic, form.BirthDate)
        if err != nil{
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        
        c.JSON(http.StatusOK, gin.H{})

    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не разрешён"})
    }
}


func (h *Handler) addPublication(c *gin.Context) {
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
        c.HTML(http.StatusOK, "addPublication.tmpl", gin.H{})
    case http.MethodPost:
        var form createPublicationForm
        if err := c.BindJSON(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
            return
        }

        var yearInt int
        if form.Year != "" {
            y, err := strconv.Atoi(form.Year)
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Год издания должен быть числом"})
                return
            }
            yearInt = y
        }

        var isbns []string
        if form.ISBN != "" {
            parts := strings.Split(form.ISBN, "|")
            for _, p := range parts {
                isbns = append(isbns, p)
            }
        }

        var otherIsbns []string
        if form.OtherIsbn != "" {
            parts := strings.Split(form.OtherIsbn, "|")
            for _, p := range parts {
                otherIsbns = append(isbns, p)
            }
        }

        var bbks []string
        if form.BBK != "" {
            parts := strings.Split(form.BBK, "+")
            for _, p := range parts {
                bbks = append(bbks, p)
            }
        }

        var otherIndexes []string
        if form.OtherIndexes != "" {
            // предположим, что они тоже разделены | (можно уточнить у фронта)
            parts := strings.Split(form.OtherIndexes, "|")
            for _, p := range parts {
                otherIndexes = append(otherIndexes, p)
            }
        }

        err := h.service.CreatePublication(
            form.Title,
            yearInt,
            form.AuthorIds,
            isbns,
            otherIsbns,
            bbks,
            otherIndexes,
        )
        if err != nil{
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }else{
            c.JSON(http.StatusOK, gin.H{})
            return
        }



    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не разрешён"})
    }
}
