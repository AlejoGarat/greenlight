package data

import (
	"time"

	"greenlight/internal/models"
)

type Movie struct {
	ID        int64          `json:"id"`
	CreatedAt time.Time      `json:"-"` // omit from output
	Title     string         `json:"title"`
	Year      int32          `json:"year"`
	Runtime   models.Runtime `json:"runtime"`
	Genres    []string       `json:"genres"`
	Version   int32          `json:"version"`
}
