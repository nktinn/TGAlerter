package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Router struct {
	db *sqlx.DB
}

func NewRouter(db *sqlx.DB) *Router {
	return &Router{db: db}
}

func (r *Router) GetRoute(service string) int64 {
	var userID int64

	query := fmt.Sprintf("SELECT userID FROM %s WHERE serviceID = $1", routes)
	row := r.db.QueryRow(query, service)
	if err := row.Scan(&userID); err != nil {
		return 0
	}

	return userID
}
