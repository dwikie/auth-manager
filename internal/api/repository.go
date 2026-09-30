package api

import (
	"database/sql"

	"github.com/dwikie/auth-manager/internal/repo"
)

type Repository struct {
	User *repo.UserRepository
}

func RegisterRepositories(db *sql.DB) Repository {

	userRepo := repo.NewUserRepository(db)

	return Repository{
		User: userRepo,
	}
}
