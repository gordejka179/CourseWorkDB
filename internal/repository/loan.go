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

//сделать выдачу, используя читательский номер, почту библиотекаря, инвентарному номеру экземпляра
func (r *Repository) MakeLoan(readerLibraryCard string, emailLibrarian string, inventoryNumber string) error {
    tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly:  false,
    })

    if err != nil {
        return fmt.Errorf("не удалось начать транзакцию: %w", err)
    }
    defer tx.Rollback()

    // Получить librarianId по email
    var librarianId int
    err = tx.QueryRow(`SELECT librarianId FROM Librarian WHERE email = $1`, emailLibrarian).Scan(&librarianId)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("библиотекарь с email %s не найден", emailLibrarian)
        }
        return fmt.Errorf("ошибка получения librarianId: %w", err)
    }


	// Получить readerId по readerLibraryCard
	var readerID int
    err = tx.QueryRow(`SELECT readerId FROM Reader WHERE libraryCard = $1`, readerLibraryCard).Scan(&readerID)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("читатель с libraryCard %s не найден", readerLibraryCard)
        }
        return fmt.Errorf("ошибка получения readerId: %w", err)
    }

    //Получить copyId по inventoryNumber
    var copyId int
    err = tx.QueryRow(`SELECT copyId FROM Copy WHERE inventoryNumber = $1`, inventoryNumber).Scan(&copyId)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("экзмемляр с inventoryNumber %s не найден", inventoryNumber)
        }
        return fmt.Errorf("ошибка получения copyId: %w", err)
    }

	var success bool
    // Выдача экземпляра с copyId библиотекарем с librarianId для читателя с readerId
	err = tx.QueryRow(`SELECT makeLoan($1, $2, $3)`, readerID, librarianId, copyId).Scan(&success)
	if err != nil {
        // Разбор ошибок:
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "ML001":
                return fmt.Errorf("экземпляр книги не найден")
            case "ML002":
                return fmt.Errorf("экземпляр книги уже кем-то забронирован(срок бронирования не истёк) или получен")
            case "ML003":
                return fmt.Errorf("У читателя на руках уже есть экземпляр этого издания")
            }
        }
        return fmt.Errorf("ошибка при бронировании: %w", err)
    }

    if !success {
        return fmt.Errorf("не удалось сделать выдачу")
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit транзакции: %w", err)
    }
    return nil
}

//используя почту читателя получить информацию о забронированных читателем экземплярах
func (r *Repository) GetLoanedBooksByReaderEmail(emailReader string) ([]models.IssueInformation, error) {
    // Получить readerId по email
    var readerID int
    err := r.db.QueryRow("SELECT readerId FROM Reader WHERE email = $1", emailReader).Scan(&readerID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return []models.IssueInformation{}, fmt.Errorf("читатель с email %s не найден", emailReader)
        }
        return []models.IssueInformation{}, fmt.Errorf("ошибка запроса: %w", err)
    }
    
    rows, err := r.db.Query(`SELECT 
    COALESCE(copyid, 0) AS copyid,
    COALESCE(expirydate::TEXT, '') AS expiryDate,
    COALESCE(inventorynumber, '') AS inventorynumber,
    COALESCE(title, '') AS title,
    COALESCE(publicationyear, 0) AS publicationyear,
    COALESCE(authors, ARRAY[]::TEXT[]) AS authors,
    COALESCE(isbns, ARRAY[]::TEXT[]) AS isbns,
    COALESCE(bbks, ARRAY[]::TEXT[]) AS bbks,
    COALESCE(otherindexes, ARRAY[]::TEXT[]) AS otherindexes,
    COALESCE(buildingid, 0) AS buildingid,
    COALESCE(buildingaddress, '') AS buildingaddress
    FROM get_current_loans_by_readerId($1)`, readerID)

    if err != nil {
        return nil, fmt.Errorf("ошибка запроса выданных книг: %w", err)
    }
    defer rows.Close()

    var result []models.IssueInformation
    for rows.Next() {
        var (
            copyId int
            expiryDate string
            inventoryNumber string
            title string
            publicationYear int
            authors []string
            isbns []string
            bbks []string
            otherIndexes []string
            buildingId int
            buildingAddress string
        )
        if err := rows.Scan(
            &copyId,
            &expiryDate,
            &inventoryNumber,
            &title,
            &publicationYear,
            pq.Array(&authors),
            pq.Array(&isbns),
            pq.Array(&bbks),
            pq.Array(&otherIndexes),
            &buildingId,
            &buildingAddress,
        ); err != nil {
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
        
        info := models.IssueInformation{
            CopyId: copyId,
            ExpiryDate: expiryDate,
            InventoryNumber: inventoryNumber,
            Title: title,
            PublicationYear: publicationYear,
            Authors: formattedAuthors,
            Isbns: isbns,
            BBKs: bbks,
            OtherIndexes: otherIndexes,
            Building: models.Building{
                BuildingId: buildingId,
                Address: buildingAddress,
            },
        }
        result = append(result, info)
    }

    return result, rows.Err()
}


//используя читательский билет читателя получить информацию о забронированных читателем экземплярах
func (r *Repository) GetLoanedBooksByReaderLibraryCard(readerLibraryCard string) ([]models.IssueInformation, error) {
    // Получить readerId по email
    var readerID int
    err := r.db.QueryRow("SELECT readerId FROM Reader WHERE librarycard = $1", readerLibraryCard).Scan(&readerLibraryCard)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return []models.IssueInformation{}, fmt.Errorf("читатель с librarycard %s не найден", readerLibraryCard)
        }
        return []models.IssueInformation{}, fmt.Errorf("ошибка запроса: %w", err)
    }
    
    rows, err := r.db.Query(`SELECT 
    COALESCE(copyid, 0) AS copyid,
    COALESCE(expirydate::TEXT, '') AS expiryDate,
    COALESCE(inventorynumber, '') AS inventorynumber,
    COALESCE(title, '') AS title,
    COALESCE(publicationyear, 0) AS publicationyear,
    COALESCE(authors, ARRAY[]::TEXT[]) AS authors,
    COALESCE(isbns, ARRAY[]::TEXT[]) AS isbns,
    COALESCE(bbks, ARRAY[]::TEXT[]) AS bbks,
    COALESCE(otherindexes, ARRAY[]::TEXT[]) AS otherindexes,
    COALESCE(buildingid, 0) AS buildingid,
    COALESCE(buildingaddress, '') AS buildingaddress
    FROM get_current_loans_by_readerId($1)`, readerID)

    if err != nil {
        return nil, fmt.Errorf("ошибка запроса выданных книг: %w", err)
    }
    defer rows.Close()

    var result []models.IssueInformation
    for rows.Next() {
        var (
            copyId int
            expiryDate string
            inventoryNumber string
            title string
            publicationYear int
            authors []string
            isbns []string
            bbks []string
            otherIndexes []string
            buildingId int
            buildingAddress string
        )
        if err := rows.Scan(
            &copyId,
            &expiryDate,
            &inventoryNumber,
            &title,
            &publicationYear,
            pq.Array(&authors),
            pq.Array(&isbns),
            pq.Array(&bbks),
            pq.Array(&otherIndexes),
            &buildingId,
            &buildingAddress,
        ); err != nil {
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
        
        info := models.IssueInformation{
            CopyId:          copyId,
            InventoryNumber: inventoryNumber,
            Title:           title,
            PublicationYear: publicationYear,
            Authors:         formattedAuthors,
            Isbns:           isbns,
            BBKs:            bbks,
            OtherIndexes:    otherIndexes,
            Building: models.Building{
                BuildingId: buildingId,
                Address:    buildingAddress,
            },
        }
        result = append(result, info)
    }

    return result, rows.Err()
}