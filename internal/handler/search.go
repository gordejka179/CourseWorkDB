package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/gordejka179/CourseWorkDB/internal/usecase"
)


type SearchForm struct {
	Authors string `json:"authors"`
	ISBN string `json:"isbn"`
	Title string `json:"title"`
	BBKs string `json:"bbks"`
	PublicationYear string `json:"publicationyear"`
	OtherIndex string `json:"otherindex"`
	AlternativeSearch bool `json:"alternativesearch"` //нужно, чтобы понимать, показывать ли ббк, которые являются рекомендательными
}


func (h *Handler) search(c *gin.Context) {
	c.HTML(http.StatusOK, "search.tmpl", gin.H{})
}

func (h *Handler) searchBook(c *gin.Context) {
	var searchForm SearchForm
    if err := c.BindJSON(&searchForm); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
        return
    }

	//Будем получить информацию об изданиях, причем будем стараться получать как можно меньше записей,
	//то есть делать поиск по параметрам, которые вернут мало значений,
	//а затем делать пересечения с другими данными


	//[Пушкин|Александр|Сергеевич];[Есенин|Сергей]
	var formattedAuthors []models.Author
	if searchForm.Authors != ""{
		Authors := strings.Split(searchForm.Authors, ";")


		for _ , a := range Authors{
			a = a[1:len(a) - 1]
			fullname := strings.Split(a, "|")

			author := models.Author{LastName: fullname[0], FirstName: fullname[1], Patronymic: fullname[2]}

			formattedAuthors = append(formattedAuthors, author)
		}
	}


	ParametersForm := usecase.ParametersForm{
		Authors: formattedAuthors,
		ISBN: searchForm.ISBN,
		Title: searchForm.Title,
		PublicationYear: searchForm.PublicationYear,
		OtherIndex: searchForm.OtherIndex,
		AdditionalSearch: searchForm.AlternativeSearch,
	}

	_, err := h.service.SearchPublications(ParametersForm)


	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка на сервере"})
        return
	}


}

