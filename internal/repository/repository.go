package repository

import (
	"github.com/jmoiron/sqlx"
)

const (
	routes = "notification_routes"
)

type RouteService interface {
	GetRoute(service string) int64
}

type Repository struct {
	RouteService
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		RouteService: NewRouter(db),
	}
}
