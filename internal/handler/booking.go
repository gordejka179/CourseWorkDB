package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//сделать бронирование
func (h *Handler) reserve(c *gin.Context) {
    emailRaw, _ := c.Get("email")

    email, ok := emailRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения email"})
        return
    }

    var reserveReq struct {
        CopyId string `json:"copyid"`
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


// для читателя: получить информацию о своих бронированиях
// для библиотекаря: по читательскому номеру хотим видеть, какие книги забронированы и возможно сделать выдачу по этой брони
func (h *Handler) getCurrentBookings(c *gin.Context) {
    roleRaw, _ := c.Get("role")
    role, ok := roleRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения role"})
        return
    }

    switch c.Request.Method {
    case http.MethodGet:
	    if role == "reader"{
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
                return
        }
        
        if role == "librarian"{        
	        c.HTML(http.StatusOK, "bookingsLibrarian.tmpl", gin.H{})
            return
        }
    case http.MethodPost:
        if role == "reader"{
		    c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        
        if role == "librarian"{

            var bookingsReq struct {
                ReaderLibraryCard string `json:"readerLibraryCard"`
            }

            if err := c.BindJSON(&bookingsReq); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
                return
            }


	        bookingsInfo, err := h.service.GetCurrentBookingsByReaderLibraryCard(bookingsReq.ReaderLibraryCard)


            if err != nil{
		        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
	        }

            result := make([]gin.H, 0, len(bookingsInfo))
            for _, ii := range bookingsInfo {
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
    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не разрешён"})
    }
    

}
