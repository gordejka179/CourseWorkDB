package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/lib/pq"
)


func (r *Repository) SearchAuthors(lastName string, firstName string, patronymic string, birthDate string) ([]models.AuthorForAdd, error) {
    query :=
     `SELECT
    authorId,
    lastName,
    COALESCE(firstName, '') AS firstName,
    COALESCE(patronymic, '') AS patronymic,
    COALESCE(birthDate::text, '{}') AS birthDate
    FROM search_authors($1, $2, $3, $4)`

	var parsed time.Time
    var rows *sql.Rows
    var err error
	if birthDate != "" {
        // Парсим в формате  год-месяц-день)
        parsed, err = time.Parse("1991-09-09", birthDate)
        if err != nil{
            return []models.AuthorForAdd{}, fmt.Errorf("неверный формат даты рождения: %w", err)
        }
		rows, err = r.db.Query(query, lastName, firstName, patronymic, parsed)

    } else {
		rows, err = r.db.Query(query, lastName, firstName, patronymic, nil)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var authors []models.AuthorForAdd
    for rows.Next() {

        var (
            authorId int
            lastName string
            firstName string
            patronymic string
            birthDate string
        )

        err = rows.Scan(&authorId, &lastName, &firstName, &patronymic, &birthDate)


        if err != nil {
            return nil, err
        }


        authors = append(authors, models.AuthorForAdd{
            AuthorId : authorId,
			FirstName: firstName,
            LastName: lastName,
            Patronymic: patronymic,
            BirthDate: birthDate,
        })

    }

    if err = rows.Err(); err != nil {
        return nil, err
    }
    return authors, nil

}

func (r *Repository) CreateAuthor(lastName, firstName, patronymic, birthDate string) error {
    tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly:  false,
    })
    if err != nil {
        return fmt.Errorf("ошибка начала транзакции: %w", err)
    }
    defer tx.Rollback()

	var parsedDate time.Time
	if birthDate != "" {
        // Парсим в формате  год-месяц-день)
        parsedDate, err = time.Parse("1991-09-09", birthDate)
        if err != nil{
            return fmt.Errorf("неверный формат даты рождения: %w", err)
        }
    }

    var patronymicPtr interface{}
    if patronymic == "" {
        patronymicPtr = nil
    } else {
        patronymicPtr = patronymic
    }

    _, err = tx.Exec(
        `SELECT create_author($1, $2, $3, $4)`,
        lastName, firstName, patronymicPtr, parsedDate,
    )
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            return fmt.Errorf("автор с такими данными уже существует")
        }
        return fmt.Errorf("ошибка вызова create_author: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit tx: %w", err)
    }
    return nil
}


func (r *Repository) CreatePublication(title string, publicationYear int, authorIds []int, isbns []string, otherIsbns []string, bbks []string, otherIndexes []string) error {
    ctx := context.Background()
    tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly:  false,
    })
    if err != nil {
        return fmt.Errorf("ошибка начала транзакции: %w", err)
    }
    defer tx.Rollback()

    var authorIdsParam interface{} = pq.Array(authorIds)
    if len(authorIds) == 0 {
        authorIdsParam = nil
    }

    var isbnsParam interface{} = pq.Array(isbns)
    if len(isbns) == 0 {
        isbnsParam = nil
    }

    var otherIsbnsParam interface{} = pq.Array(otherIsbns)
    if len(otherIsbns) == 0 {
        otherIsbnsParam = nil
    }

    var bbksParam interface{} = pq.Array(bbks)
    if len(bbks) == 0 {
        bbksParam = nil
    }

    var otherIndexesParam interface{} = pq.Array(otherIndexes)
    if len(otherIndexes) == 0 {
        otherIndexesParam = nil
    }

    // Вызов хранимой функции
    query := `SELECT create_publication($1, $2, $3, $4, $5, $6, $7)`
    var publicationId int
    err = tx.QueryRowContext(ctx, query,
        title,
        publicationYear,
        authorIdsParam,
        isbnsParam,
        otherIsbnsParam,
        bbksParam,
        otherIndexesParam,
    ).Scan(&publicationId)

    if err != nil {
        // Разбор ошибок:
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "A0001":
                return fmt.Errorf("Таких авторов нет в базе")
            case "I0001":
                return fmt.Errorf("Такой isbn уже есть у издания в базе")
            }
        }
        return fmt.Errorf("ошибка при вызове create_publication: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("ошибка коммита: %w", err)
    }

    return nil
}
