package data

import (
	"database/sql/driver"
	"time"

	"greenlight/internal/models"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64          `json:"id" db:"id"`
	CreatedAt time.Time      `json:"-" db:"created_at"` // omit from output
	Title     string         `json:"title" db:"title"`
	Year      int32          `json:"year" db:"year"`
	Runtime   models.Runtime `json:"runtime" db:"runtime"`
	Genres    pqStrSlice     `json:"genres" db:"genres"`
	Version   int32          `json:"version" db:"version"`
}

type pqStrSlice []string

func (p *pqStrSlice) Scan(src interface{}) error {
	arr := pq.StringArray(*p)
	return arr.Scan(src)
}

func (p pqStrSlice) Value() (driver.Value, error) {
	return pq.Array(p).Value()
}
