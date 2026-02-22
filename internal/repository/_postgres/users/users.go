package users

import (
	"fmt"
	"time"

	_postgres "golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT id, name, email, age, phone, created_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %w", err)
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "SELECT id, name, email, age, phone, created_at FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: user with id=%d not found", id)
	}
	return &user, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	var newID int
	query := `INSERT INTO users (name, email, age, phone) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.DB.QueryRow(query, user.Name, user.Email, user.Age, user.Phone).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}
	return newID, nil
}

func (r *Repository) UpdateUser(id int, user modules.User) error {
	query := `UPDATE users SET name = $1, email = $2, age = $3, phone = $4 WHERE id = $5`
	result, err := r.db.DB.Exec(query, user.Name, user.Email, user.Age, user.Phone, id)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateUser RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("UpdateUser: user with id=%d does not exist", id)
	}
	return nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return 0, fmt.Errorf("DeleteUser: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("DeleteUser RowsAffected: %w", err)
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("DeleteUser: user with id=%d does not exist", id)
	}
	return rowsAffected, nil
}
