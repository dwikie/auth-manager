package repo

import (
	"database/sql"
)

type Repository struct {
	User *UserRepository
}

func Initialize(db *sql.DB) *Repository {
	userRepo := NewUserRepository(db)

	return &Repository{
		User: userRepo,
	}
}
