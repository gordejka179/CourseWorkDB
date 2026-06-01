package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/lib/pq"
)


func (r *Repository) SearchAuthors(lastName string, firstName string, patronymic string, birthYear string) ([]models.AuthorForAdd, error) {
    query :=
     `SELECT
    authorId,
    lastName,
    COALESCE(firstName, '') AS firstName,
    COALESCE(patronymic, '') AS patronymic,
    COALESCE(birthYear::text, '') AS birthYear
    FROM search_authors($1, $2, $3, $4)`


    var parsedYear interface{} = nil
    var err error

    var rows *sql.Rows
	if birthYear != "" {
        parsedYear, err = strconv.Atoi(birthYear)
        if err != nil{
            return []models.AuthorForAdd{}, fmt.Errorf("неверный формат года рождения: %w", err)
        }
		rows, err = r.db.Query(query, lastName, firstName, patronymic, parsedYear)

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
            birthYear string
        )

        err = rows.Scan(&authorId, &lastName, &firstName, &patronymic, &birthYear)


        if err != nil {
            return nil, err
        }


        authors = append(authors, models.AuthorForAdd{
            AuthorId: authorId,
			FirstName: firstName,
            LastName: lastName,
            Patronymic: patronymic,
            BirthYear: birthYear,
        })

    }

    if err = rows.Err(); err != nil {
        return nil, err
    }
    return authors, nil

}

func (r *Repository) CreateAuthor(lastName, firstName, patronymic, birthYear string) error {
    tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly:  false,
    })
    if err != nil {
        return fmt.Errorf("ошибка начала транзакции: %w", err)
    }
    defer tx.Rollback()

    var parsedYear interface{} = nil
	if birthYear != "" {
        parsedYear, err = strconv.Atoi(birthYear)
        if err != nil{
            return fmt.Errorf("неверный формат года рождения: %w", err)
        }
    }

    var firstNamePtr interface{}
    if firstName == "" {
        firstNamePtr = nil
    } else {
        firstNamePtr = firstName
    }

    var patronymicPtr interface{}
    if patronymic == "" {
        patronymicPtr = nil
    } else {
        patronymicPtr = patronymic
    }

    _, err = tx.Exec(
        `SELECT create_author($1, $2, $3, $4)`,
        firstNamePtr, lastName, patronymicPtr, parsedYear,
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
