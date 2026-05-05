package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//делаем выдачу экземпляра
func (h *Handler) makeLoan(c *gin.Context) {
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

    emailLibrarianRaw, _ := c.Get("email")

    emailLibrarian, ok := emailLibrarianRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения emailLibrarian"})
        return
    }

    switch c.Request.Method {
    case http.MethodPost:
        var loanReq struct {
            ReaderLibraryCard string `json:"readerlibrarycard"`
            InventoryNumber string `json:"inventorynumber"`
        }

        if err := c.BindJSON(&loanReq); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
            return
        }

        fmt.Println(loanReq.ReaderLibraryCard, emailLibrarian, loanReq.InventoryNumber)
        err := h.service.MakeLoan(loanReq.ReaderLibraryCard, emailLibrarian, loanReq.InventoryNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Книга забронирована"})

    case http.MethodGet:
        c.HTML(http.StatusOK, "loansLibrarian.tmpl", gin.H{})
    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не разрешён"})
    }
}


//получение
// для читателя: списка выданных ему экземпляров
// для библиотекаря: по читательскому билету читателя книги, которые у него на руках
func (h *Handler) getLoanedBooks(c *gin.Context) {
    roleRaw, _ := c.Get("role")
    role, ok := roleRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения role"})
        return
    }

    switch c.Request.Method {
    case http.MethodGet:
	    if role == "librarian"{
		    c.HTML(http.StatusOK, "loansLibrarian.tmpl", gin.H{})
            return
        }
        
        if role == "reader"{        
	        emailReaderRaw, _ := c.Get("email")
            emailReader, ok := emailReaderRaw.(string)
            if !ok {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения emailReader"})
                return
            }

	        IssueInformationArray, err := h.service.GetLoanedBooksByReaderEmail(emailReader)

            if err != nil{
		        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
	        }

            result := make([]gin.H, 0, len(IssueInformationArray))
            for _, ii := range IssueInformationArray {
		        authorStrings := make([]string, len(ii.Authors))

		        //из массива структур Author делаем массив строк
		        for i, a := range ii.Authors {
    		        authorStrings[i] = strings.TrimSpace(fmt.Sprintf("%s %s %s", a.LastName, a.FirstName, a.Patronymic))
		        }

                item := gin.H{
                    "copyid": ii.CopyId,
                    "expirydate": ii.ExpiryDate,
                    "inventorynumber": ii.InventoryNumber,
                    "title": ii.Title,
                    "publicationyear": ii.PublicationYear,
                    "authors": authorStrings,
                    "isbns": ii.Isbns,
	                "bbks": ii.BBKs,
                    "otherindexes": ii.OtherIndexes,
                    "buildingaddress": ii.Building.Address,
                }
                result = append(result, item)
            }

            c.HTML(http.StatusOK, "loansReader.tmpl", result)
            return
        }
    case http.MethodPost:
        if role == "reader"{
		    c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        
        if role == "librarian"{        
	        librarycardRaw, _ := c.Get("librarycard")
            librarycard, ok := librarycardRaw.(string)
            if !ok {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения librarycard"})
                return
            }

	        IssueInformationArray, err := h.service.GetLoanedBooksByReaderLibraryCard(librarycard)

            if err != nil{
		        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
	        }

            result := make([]gin.H, 0, len(IssueInformationArray))
            for _, ii := range IssueInformationArray {
		        authorStrings := make([]string, len(ii.Authors))

		        //из массива структур Author делаем массив строк
		        for i, a := range ii.Authors {
    		        authorStrings[i] = strings.TrimSpace(fmt.Sprintf("%s %s %s", a.LastName, a.FirstName, a.Patronymic))
		        }

                item := gin.H{
                    "copyid": ii.CopyId,
                    "expirydate": ii.ExpiryDate,
                    "inventorynumber": ii.InventoryNumber,
                    "title": ii.Title,
                    "publicationyear": ii.PublicationYear,
                    "authors": authorStrings,
                    "isbns": ii.Isbns,
	                "bbks": ii.BBKs,
                    "otherindexes": ii.OtherIndexes,
                    "buildingaddress": ii.Building.Address,
                }
                result = append(result, item)
            }

            c.JSON(http.StatusOK, result)
            return
        }
    }
}
