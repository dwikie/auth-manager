package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dwikie/auth-manager/internal/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User, profile *domain.UserProfile) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var userId int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (username, password, status) 
		VALUES ($1, $2, $3) RETURNING id`,
		user.Username, user.Password, 0,
	).Scan(&userId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO user_profile (user_id, first_name, last_name, email, phone_number) 
		VALUES ($1, $2, $3, $4, $5)`,
		userId, profile.FirstName, profile.LastName, profile.Email, profile.Phone,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT id, username, password, status, created_at, updated_at 
		FROM users 
		WHERE username = $1`,
		username,
	)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserProfileByUserId(ctx context.Context, userId int) (*domain.UserProfile, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT id, user_id, first_name, last_name, email, email_status, phone, phone_status 
		FROM user_profiles 
		WHERE user_id = $1`,
		userId,
	)

	var profile domain.UserProfile
	err := row.Scan(&profile.ID, &profile.UserId, &profile.FirstName, &profile.LastName, &profile.Email, &profile.EmailStatus, &profile.Phone, &profile.PhoneStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user profile not found") // User profile not found
		}
		return nil, err
	}

	return &profile, nil
}
