package repository

import (
	"strings"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	"github.com/lib/pq"
)

//получаем издание(издания, если есть записи в otherISBN) по isbn
func (r *Repository) GetPublicationsByISBN(ISBN string) ([]models.Publication, error) {
    query := `SELECT 
    publicationId, 
    COALESCE(title, '') as title,
    COALESCE(publicationYear, 0) AS publicationYear,
    COALESCE(isbns, '{}') AS isbns,
    COALESCE(bbks, '{}') AS bbks,
    COALESCE(otherIndexes, '{}') AS otherIndexes,
    COALESCE(authors, '{}') AS authors
    FROM search_publications_by_isbn($1)`

    rows, err := r.db.Query(query, ISBN)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var publications []models.Publication

    for rows.Next() {
        var (
            publicationId int
            title string
            publicationYear int
            bbks []string
            otherIndexes []string
            isbns []string
            authors []string //ниже преобразуем
        )

        err := rows.Scan(
            &publicationId,
            &title,
            &publicationYear,
            pq.Array(&isbns),
            pq.Array(&bbks),
            pq.Array(&otherIndexes),
            pq.Array(&authors),
        )
        
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

	    

        pub := models.Publication{
            ID: publicationId,
            Title: title,
            PublicationYear: publicationYear,
            Authors: formattedAuthors,
            ISBNs: isbns,
            BBKs: bbks,
        }
        publications = append(publications, pub)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return publications, nil
}


func (r *Repository) GetCopiesByIDList(ids []int) ([]models.Copy, error) {
   query := 
     `SELECT 
    copyId,
    publicationId,
    buildingId,
    COALESCE(readerId, 0) AS readerId,
    COALESCE(librarianId, 0) AS librarianId,
    COALESCE(address, '') AS address,
    COALESCE(description, '') AS description
    FROM get_copies_info_by_ids($1)`

    rows, err := r.db.Query(query, pq.Array(ids))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var copies []models.Copy
    for rows.Next() {
        var c models.Copy
        err := rows.Scan(
            &c.CopyId,
            &c.PublicationId,
            &c.BuildingId,
            &c.ReaderId,
            &c.LibrarianId,
            &c.Address,
            &c.Description,
        )
        if err != nil {
            return nil, err
        }
        copies = append(copies, c)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }
    return copies, nil
}


//получаем издания по названию
func (r *Repository) GetPublicationsByTitle(title string) ([]models.Publication, error){
    query := `SELECT 
    publicationId, 
    COALESCE(title, '') as title,
    COALESCE(publicationYear, 0) AS publicationYear,
    COALESCE(isbns, '{}') AS isbns,
    COALESCE(bbks, '{}') AS bbks,
    COALESCE(otherIndexes, '{}') AS otherIndexes,
    COALESCE(authors, '{}') AS authors
    FROM search_publications_by_title($1)`

    rows, err := r.db.Query(query, title)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var publications []models.Publication

    for rows.Next() {
        var (
            publicationId int
            title string
            publicationYear int
            bbks []string
            otherIndexes []string
            isbns []string
            authors []string //ниже преобразуем
        )

        err := rows.Scan(
            &publicationId,
            &title,
            &publicationYear,
            pq.Array(&isbns),
            pq.Array(&bbks),
            pq.Array(&otherIndexes),
            pq.Array(&authors),
        )

        
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
	    

        pub := models.Publication{
            ID: publicationId,
            Title: title,
            PublicationYear: publicationYear,
            Authors: formattedAuthors,
            ISBNs: isbns,
            BBKs: bbks,
        }
        publications = append(publications, pub)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return publications, nil
}
