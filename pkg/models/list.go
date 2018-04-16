package models

import (
	"time"
)

// List contains information about a list
type List struct {
	ID      int64     `db:"ID"`
	Name    string    `db:"Name" binding:"required,min=1,max=45"`
	Comment string    `db:"Comment" binding:"max=255"`
	Created time.Time `db:"Created"`
	Updated time.Time `db:"Updated"`
}
