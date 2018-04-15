package models

import (
	"time"
)

// List contains information about a list
type List struct {
	ID      int       `db:"ID"`
	Name    string    `db:"Name"`
	Comment string    `db:"Comment"`
	Created time.Time `db:"Created"`
	Updated time.Time `db:"Updated"`
}
