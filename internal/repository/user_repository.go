package repository

import (
	"database/sql"

	"myapp/internal/domain"
)

type UserRepository interface {
	GetByID(id int64) (*domain.User, error)
	Create(user *domain.User) error
	GetAll() ([]*domain.User, error)
	Update(user *domain.User) error
	Delete(id int64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, email FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.Email).Scan(&user.ID)
}

func (r *userRepository) GetAll() ([]*domain.User, error) {
	query := `SELECT id, username, email FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`
	result, err := r.db.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
