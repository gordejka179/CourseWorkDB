package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gordejka179/CourseWorkDB/internal/models"
)


func (r *Repository) GetOverdueCopies() ([]models.OverdueCopy, error) {
    query := `SELECT 
        copy_id,
        inventory_number,
        book_title,
        expiry_date,
        days_overdue,
        reader_last_name,
        reader_first_name,
        reader_patronymic,
        reader_email,
        reader_library_card
    FROM get_all_overdue_copies()`

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("ошибка выполнения запроса get_all_overdue_copies: %w", err)
    }
    defer rows.Close()

    var overdueList []models.OverdueCopy
    for rows.Next() {
        var oc models.OverdueCopy
        var patronymic sql.NullString

        err := rows.Scan(
            &oc.CopyID,
            &oc.InventoryNumber,
            &oc.BookTitle,
            &oc.ExpiryDate,
            &oc.DaysOverdue,
            &oc.ReaderLastName,
            &oc.ReaderFirstName,
            &patronymic,
            &oc.ReaderEmail,
            &oc.ReaderLibraryCard,
        )
        if err != nil {
            return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
        }
        if patronymic.Valid {
            oc.ReaderPatronymic = patronymic.String
        }
        overdueList = append(overdueList, oc)
    }
    return overdueList, rows.Err()
}


func (r *Repository) GetOverallStats() (*models.OverallStats, error) {
    var stats models.OverallStats
    row := r.db.QueryRow(`SELECT report_overall()`)
    var jsonData []byte
    if err := row.Scan(&jsonData); err != nil {
        return nil, fmt.Errorf("query overall stats: %w", err)
    }
    if err := json.Unmarshal(jsonData, &stats); err != nil {
        return nil, fmt.Errorf("unmarshal overall stats: %w", err)
    }
    return &stats, nil
}