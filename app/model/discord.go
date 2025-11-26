package model

import (
	"database/sql"
)

type Discord struct {
	db *sql.DB
}

func NewDiscord(db *sql.DB) *Discord {
	discord := Discord{db: db}

	return &discord
}
