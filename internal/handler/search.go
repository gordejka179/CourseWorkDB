package handler

import (
	"fmt"
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


//выдаём html
func (h *Handler) search(c *gin.Context) {
    roleRaw, _ := c.Get("role")
    role, ok := roleRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения role"})
        return
    }

	if role == "reader"{
        c.HTML(http.StatusOK, "searchReader.tmpl", gin.H{})
        return
    }
        
    if role == "librarian"{        
	    c.HTML(http.StatusOK, "searchLibrarian.tmpl", gin.H{})
        return
    }
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

	pubsResponse, err := h.service.SearchPublications(ParametersForm)


	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка на сервере"})
        return
	}

    result := make([]gin.H, 0, len(pubsResponse))
    for _, pub := range pubsResponse {
        buildingsSlice := make([]gin.H, 0, len(pub.Buildings))

		//меняем map из зданий на слайс
        for _, b := range pub.Buildings {
            buildingsSlice = append(buildingsSlice, gin.H{
                "buildingId": b.BuildingId,
                "address": b.Address,
                "description": b.Description,
                "totalCopies": b.TotalCopies,
                "availableCopies": b.AvailableCopies,
                "availableCopyIds": b.AvailableCopyIds,
            })
        }


		authorStrings := make([]string, len(pub.Authors))

		//из массива структур Author делаем массив строк
		for i, a := range pub.Authors {
    		authorStrings[i] = strings.TrimSpace(fmt.Sprintf("%s %s %s", a.LastName, a.FirstName, a.Patronymic))
		}

        item := gin.H{
            "id": pub.Id,
            "title": pub.Title,
            "publicationyear": pub.PublicationYear,
            "authors": authorStrings,
            "isbns": pub.Isbns,
            "bbks": pub.BBKs,
            "otherindexes": pub.OtherIndexes,
            "buildings": buildingsSlice,
        }
        result = append(result, item)
    }

	c.JSON(http.StatusOK, result)

}

