package repository

import (
	"database/sql"
	"fmt"

	"github.com/gordejka179/CourseWorkDB/internal/models"
)

func (r *Repository) ReturnBook(inventoryNumber string) error {
    tx, err := r.db.Begin()
    if err != nil {
        return fmt.Errorf("не удалось начать транзакцию: %w", err)
    }
    defer tx.Rollback()

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
        SELECT return_copy($1)`,
         copyId).Scan(&success)
    if err != nil {
        return fmt.Errorf("ошибка вызова return_copy: %w", err)
    }

    if !success {
        return fmt.Errorf("экземпляр с инвентарным номером %s не найден", inventoryNumber)
    }

    //фиксируем транзакцию
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("ошибка фиксации транзакции: %w", err)
    }

    return nil
}

func (r *Repository) GetBuildings()([]models.Building2, error) {
    query := `SELECT libraryBuildingId, address, description FROM LibraryBuilding ORDER BY libraryBuildingId`

    rows, err := r.db.Query(query)
    if err != nil {
        fmt.Println(err)
        return nil, fmt.Errorf("ошибка выполнения запроса GetBuildings: %w", err)
    }
    defer rows.Close()

    var buildings []models.Building2
    for rows.Next() {
        var b models.Building2
        err := rows.Scan(&b.Id, &b.Address, &b.Description)
        if err != nil {
            return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
        }
        buildings = append(buildings, b)
    }
    return buildings, rows.Err()
}
