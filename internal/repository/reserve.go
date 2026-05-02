package repository

import (
	"context"
	"database/sql"
	"fmt"

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