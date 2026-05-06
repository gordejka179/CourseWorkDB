package repository

import (
	"database/sql"
	"fmt"
)

func (r *Repository) ReturnBook(readerLibraryCard string, inventoryNumber string) error {
    tx, err := r.db.Begin()
    if err != nil {
        return fmt.Errorf("не удалось начать транзакцию: %w", err)
    }
    defer tx.Rollback()

    //reader_id по номеру читательского билета
    var readerId int
    err = tx.QueryRow(`
        SELECT readerId 
        FROM Reader 
        WHERE libraryCard = $1
    `, readerLibraryCard).Scan(&readerId)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("читатель с номером %s не найден", readerLibraryCard)
        }
        return fmt.Errorf("ошибка поиска читателя: %w", err)
    }

    // copy_id по инвентарному номеру
    var copyId int
    err = tx.QueryRow(`
        SELECT copyId 
        FROM Copy 
        WHERE inventoryNumber = $1
    `, inventoryNumber).Scan(&copyId)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("экземпляр с инвентарным номером %s не найден", inventoryNumber)
        }
        return fmt.Errorf("ошибка поиска экземпляра: %w", err)
    }

    var success bool
    err = tx.QueryRow(`
        SELECT return_copy($1, $2)
    `, copyId, readerId).Scan(&success)
    if err != nil {
        return fmt.Errorf("ошибка вызова return_copy: %w", err)
    }

    if !success {
        return fmt.Errorf("экземпляр %s не числится за читателем %s", inventoryNumber, readerLibraryCard)
    }

    //фиксируем транзакцию
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("ошибка фиксации транзакции: %w", err)
    }

    return nil
}
