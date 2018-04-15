package models

import (
	"time"
)

// IP contains IP data
type IP struct {
	ID      int       `db:"ID"`
	ListID  int       `db:"ListID"`
	Address string    `db:"Address"` // https://dev.mysql.com/doc/refman/5.6/en/miscellaneous-functions.html#function_inet6-aton
	Network int8      `db:"Network"`
	Comment string    `db:"Comment"`
	Created time.Time `db:"Created"`
	Updated time.Time `db:"Updated"`
}
