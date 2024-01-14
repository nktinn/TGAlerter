package repository

type Alerter interface {
}

type Repository struct {
	Alerter
}

func NewRepository() *Repository {
	return &Repository{}
}
