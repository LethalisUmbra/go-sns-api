package models

import "time"

type Error struct {
	ID      int       `json:"id"`
	Status  int       `json:"status"`
	Error   string    `json:"error"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}
