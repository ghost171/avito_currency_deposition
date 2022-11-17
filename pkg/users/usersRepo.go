package users

import (
	"database/sql",
	"fmt"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Deposit(userID string, value float64) error {

}