package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/lib/pq"
)

func (r *Repository) GetCurrentBookingsByEmail(email string) ([]models.BookingInformation, error) {
    // Получить readerId по email
    var readerID int
    err := r.db.QueryRow("SELECT readerId FROM Reader WHERE email = $1", email).Scan(&readerID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return []models.BookingInformation{}, fmt.Errorf("читатель с email %s не найден", email)
        }
        return []models.BookingInformation{}, fmt.Errorf("ошибка запроса: %w", err)
    }
    
    query := 
     `SELECT 
    copyId,
    inventoryNumber,
    COALESCE(title, '') AS title,
    COALESCE(publicationYear, 0) AS publicationYear,
    COALESCE(authors, '{}') AS authors,
    COALESCE(isbns, '{}') AS isbn,
    COALESCE(bbks, '{}') AS bbks,
    COALESCE(otherIndexes, '{}') AS otherIndexes,
    buildingId,
    buildingAddress
    FROM get_current_bookings_by_readerId($1)`

    rows, err := r.db.Query(query, readerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var bookings []models.BookingInformation

    for rows.Next() {

        var (
            copyId int
            inventoryNumber string
            title string
            publicationYear int
            authors pq.StringArray
            isbns pq.StringArray
            bbks pq.StringArray
            otherIndexes pq.StringArray
            buildingId int
            buildingAddr string
        )
       
        err = rows.Scan(&copyId, &inventoryNumber, &title, &publicationYear, &authors, &isbns,
                    &bbks, &otherIndexes, &buildingId, &buildingAddr)


        if err != nil {
            return nil, err
        }

        var formattedAuthors []models.Author

        for _ , a := range authors{
            if a != ""{
			    fullname := strings.Split(a, "|")

			    author := models.Author{LastName: fullname[0], FirstName: fullname[1], Patronymic: fullname[2]}

			    formattedAuthors = append(formattedAuthors, author)
            }
		}

        bookings = append(bookings, models.BookingInformation{
            CopyId: copyId,
            InventoryNumber: inventoryNumber,
            Title: title,
            PublicationYear: publicationYear,
            Authors: formattedAuthors,
            Isbns: isbns,
            BBKs: bbks,
            OtherIndexes: otherIndexes,
            Building: models.Building{
                BuildingId: buildingId,
                Address:    buildingAddr,
            },
        })

    }

    if err = rows.Err(); err != nil {
        return nil, err
    }
    return bookings, nil
}

