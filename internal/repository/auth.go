package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/lib/pq"
)

func (r *Repository) CheckIfReaderExists(email string) (bool, error) {
    var exists bool
    
    query := "SELECT EXISTS(SELECT 1 FROM Reader WHERE email = $1)"
    
    err := r.db.QueryRow(query, email).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("check user exists error: %w", err)
    }
    
    return exists, nil
}


func (r *Repository) CreateReader(reader *models.Reader) error {
    err := r.createReaderTx(reader)
    if err != nil {
        return fmt.Errorf("Не получилось создать читателя: %w", err)
    }
    return nil
}


// Транзакция для создания читателя
func (r *Repository) createReaderTx(reader *models.Reader) error {
    // Устанаваливаем уровень SERIALIZABLE
    tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly:  false,
    })
    if err != nil {
        return fmt.Errorf("Ошибка создания читателя в транзакции: %w", err)
    }
    defer tx.Rollback() // откат, если не закоммитим

    // Вызываем SQL-функцию createReader, которая возвращает новый id читателя и номер читательского билета
    var newID int
    var newCard string
    err = tx.QueryRow(
        `SELECT readerId, libraryCard FROM createReader($1, $2, $3, $4, $5, $6, $7)`,
        reader.Email,
        reader.PasswordHash,
        reader.FirstName,
        reader.LastName,
        reader.PassportSeries,
        reader.PassportNumber,
        reader.Patronymic,
    ).Scan(&newID, &newCard)


     if err != nil {
        // Попытка привести ошибку к pq.Error
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "EM001":
                return fmt.Errorf("Пользователь с такой почтой уже есть")
            case "PS001":
                return fmt.Errorf("Пользователь с такими паспортными данными уже есть")
            }
        }
        return fmt.Errorf("Ошибка создания читателя в транзакции: %w", err)
    }

    reader.ReaderId = newID
    reader.LibraryCard = newCard

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit tx: %w", err)
    }
    return nil
}



// Проверка credentials читателя
func (r *Repository) CheckReaderCredentials(email, password string) (success bool, err error) {
    var isValid bool

    err = r.db.QueryRow(
        `SELECT checkReaderCredentials($1, $2)`,
        email, password,
    ).Scan(&isValid)
    fmt.Println(isValid, err)

    if err != nil {
        return false, fmt.Errorf("Ошибка при проверке данных читателя: %w", err)
    }

    return isValid, nil
}


// Проверка credentials библиотекаря
func (r *Repository) CheckLibrarianCredentials(email, password string) (success bool, err error) {
    var isValid bool

    err = r.db.QueryRow(
        `SELECT checkLibrarianCredentials($1, $2)`,
        email, password,
    ).Scan(&isValid)

    if err != nil {
        return false, fmt.Errorf("Ошибка при проверке данных библиотекаря: %w", err)
    }

    return isValid, nil
}
