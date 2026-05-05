package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/lib/pq"
)

func (r *Repository) ReserveCopyByEmail(email string, copyID int) error {
    tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly:  false,
    })

    if err != nil {
        return fmt.Errorf("не удалось начать транзакцию: %w", err)
    }
    defer tx.Rollback()

    // Получить readerId по email
    var readerID int
    err = tx.QueryRow(`SELECT readerId FROM Reader WHERE email = $1`, email).Scan(&readerID)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("читатель с email %s не найден", email)
        }
        return fmt.Errorf("ошибка получения readerId: %w", err)
    }

    // бронирование для читателя с readerId
    var success bool
    err = tx.QueryRow(`SELECT reserveCopyByEmail($1, $2)`, readerID, copyID).Scan(&success)
    if err != nil {
        // Разбор ошибок:
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "BK001":
                return fmt.Errorf("экземпляр книги не найден")
            case "BK002":
                return fmt.Errorf("экземпляр книги уже кем-то забронирован или получен")
            case "BK003":
                return fmt.Errorf("Вы уже либо держите бронь на экземпляр этого издания или имеете экземпляр на руках")
            }
        }
        return fmt.Errorf("ошибка при бронировании: %w", err)
    }

    if !success {
        return fmt.Errorf("не удалось забронировать экземпляр")
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit транзакции: %w", err)
    }
    return nil
}


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
    buildingAddress,
    COALESCE(expiryDate, '') AS expiryDate
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
            expiryDate string
        )
       
        err = rows.Scan(&copyId, &inventoryNumber, &title, &publicationYear, &authors, &isbns,
                    &bbks, &otherIndexes, &buildingId, &buildingAddr, &expiryDate)


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
            ExpiryDate: expiryDate,
        })

    }

    if err = rows.Err(); err != nil {
        return nil, err
    }
    return bookings, nil
}


func (r *Repository) GetCurrentBookingsByReaderLibraryCard(readerLibraryCard string) ([]models.BookingInformation, error) {
    // Получить readerId по readerLibraryCard
    var readerID int
    err := r.db.QueryRow("SELECT readerId FROM Reader WHERE libraryCard = $1", readerLibraryCard).Scan(&readerID)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return []models.BookingInformation{}, fmt.Errorf("читатель с libraryCard %s не найден", readerLibraryCard)
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
    buildingAddress,
    COALESCE(expiryDate, '') AS expiryDate
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
            expiryDate string
        )
       
        err = rows.Scan(&copyId, &inventoryNumber, &title, &publicationYear, &authors, &isbns,
                    &bbks, &otherIndexes, &buildingId, &buildingAddr, &expiryDate)


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
            ExpiryDate: expiryDate,
        })

    }

    if err = rows.Err(); err != nil {
        return nil, err
    }
    return bookings, nil
}
