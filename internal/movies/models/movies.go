package data

import (
	"database/sql/driver"
	"time"

	"greenlight/internal/models"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64          `json:"id"`
	CreatedAt time.Time      `json:"-" db:"created_at"` // omit from output
	Title     string         `json:"title"`
	Year      int32          `json:"year"`
	Runtime   models.Runtime `json:"runtime"`
	Genres    pqStrSlice     `json:"genres"`
	Version   int32          `json:"version"`
}

type pqStrSlice []string

func (p *pqStrSlice) Scan(src interface{}) error {
	arr := pq.StringArray(*p)
	return arr.Scan(src)
}

func (p pqStrSlice) Value() (driver.Value, error) {
	return pq.Array(p).Value()
}
