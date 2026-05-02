package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/lib/pq"
)

//сделать выдачу
func (r *Repository) MakeLoan(readerLibraryCard string, emailLibrarian string, copyId int) error {
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

//получить информацию о забронированных читателем экземплярах
func (r *Repository) GetLoanedBooks(emailReader string) ([]models.IssueInformation, error) {
	return []models.IssueInformation{}, nil
}