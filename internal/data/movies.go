package data

import "time"

type Movie struct {
	ID int64	`json:"id"`
	CreatedAt time.Time `json:"-"`
	Title string 	`json:"title"`
	Year int32 		`json:"year,omitempty"`
	Runtime Runtime	`json:"runtime,omitempty"`
	Genres []string `json:"genres,omitempty"`
	Version int32	`json:"version"`
}

type Tasks struct {
	ID int64 	`json:"id"`
	Content string `json:"content"`
	Description string `json:"description,omitempty"`
	CreatedAt time.Time	`json:"-"`
	Version int32 	`json:"version"`
}
