package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"greenlight.Arkitecth.github/internal/validator"
)


type Movie struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title string 	`json:"title"`
	Year int32 		`json:"year,omitempty"`
	Runtime Runtime	`json:"runtime,omitempty"`
	Genres []string	`json:"genres,omitempty"`
	Version int32	`json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO movies (title, year, runtime, genres) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id, created_at, version`
	
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1  {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1
	`

	var movie Movie 

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID, 
		&movie.CreatedAt, 
		&movie.Title, 
		&movie.Year,
		&movie.Runtime, 
		pq.Array(&movie.Genres), 
		&movie.Version,
	)

	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default: 
				return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
		UPDATE MOVIES
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 
		RETURNING version
	`

	args := []any{
		movie.Title, 
		movie.Year, 
		movie.Runtime, 
		pq.Array(movie.Genres), 
		movie.ID,
	}
	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	return nil
}


func ValidateMovie(v * validator.Validator, m *Movie) {
	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) <= 500, "title", "must not be more than 500 bytes long")

    v.Check(m.Year != 0, "year", "must be provided")
    v.Check(m.Year >= 1888, "year", "must be greater than 1888")
    v.Check(m.Year <= int32(time.Now().Year()), "year", "must not be in the future")

    v.Check(m.Runtime != 0, "runtime", "must be provided")
    v.Check(m.Runtime > 0, "runtime", "must be a positive integer")

    v.Check(m.Genres != nil, "genres", "must be provided")
    v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
    v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")
    // Note that we're using the Unique helper in the line below to check that all 
    // values in the m.Genres slice are unique.
    v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")

}
