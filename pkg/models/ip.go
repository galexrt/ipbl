package models

import (
	"time"
)

// IP contains IP data
type IP struct {
	ID      int64     `db:"ID"`
	ListID  int64     `db:"ListID"`
	Address string    `db:"Address" binding:"required,min=1,max=40"`
	Network int16     `db:"Network"`
	Comment string    `db:"Comment"`
	Created time.Time `db:"Created"`
	Updated time.Time `db:"Updated"`
}
