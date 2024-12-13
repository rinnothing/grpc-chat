package model

import "time"

type Message struct {
	ID     int
	UserID int
	Time   time.Time
}
