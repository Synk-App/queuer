package model

import (
	"database/sql"
)

type Telegram struct {
	db *sql.DB
}

func NewTelegram(db *sql.DB) *Telegram {
	telegram := Telegram{db: db}

	return &telegram
}
