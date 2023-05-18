package models

import (
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
	Genres    pq.StringArray `json:"genres" db:"genres"`
	Version   int32          `json:"version" db:"version"`
}
