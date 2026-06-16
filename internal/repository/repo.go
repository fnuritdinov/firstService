package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fnuritdinov/firstService/internal/models"
	errs "github.com/fnuritdinov/firstService/pkg/errors"
)

type IUserRepo interface {
	Login(ctx context.Context, auth models.Auth) error
	GetAll(ctx context.Context) ([]models.User, error)
	Get(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, id int, updatedUser string) error
	Delete(ctx context.Context, id int) error
	UpdatePassword(ctx context.Context, pass models.User) error
	GetUser(ctx context.Context, id int, auth models.User) (models.User, error)
	UpdateAge(ctx context.Context, req int, id int) error
}
type repo struct {
	db *sql.DB
}

func New(db *sql.DB) IUserRepo {
	return &repo{
		db: db,
	}
}

func (r *repo) Login(ctx context.Context, auth models.Auth) error {
	const query = `
		INSERT INTO auth (login, password, user_id) 
		VALUES ($1, $2, $3);`

	result, err := r.db.Exec(query, auth.Login, auth.Password, auth.UserID)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error from RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil

}

func (r *repo) Create(ctx context.Context, user models.User) error {

	const query = `
		INSERT INTO users (name, age, is_active) 
		VALUES ($1, $2, $3);`

	rows, err := r.db.Exec(query, user.Name, user.Age, user.IsActive)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("sql.Result.RowsAffected(): %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repo) Update(ctx context.Context, id int, updatedUser string) error {

	const query = `
		UPDATE users 
			SET name = $1 
		WHERE id = $2;`

	rows, err := r.db.Exec(query, updatedUser, id)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", rows)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("error from rowsAffected %w", err)
	}
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]models.User, error) {

	const query = `
		SELECT name, age, is_active 
			FROM users;`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.User{}, fmt.Errorf("error from r.db.Query %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.Name, &user.Age, &user.IsActive); err != nil {
			return []models.User{}, fmt.Errorf("error from rows.Scan %w", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return []models.User{}, fmt.Errorf("error from rows.Err() %w", err)
	}
	return users, nil
}

func (r *repo) Get(ctx context.Context) ([]models.User, error) {

	const query = `
		SELECT id, name, age, is_active 
			FROM users 
		WHERE is_active = $1;`

	rows, err := r.db.Query(query, true)
	if err != nil {
		return []models.User{}, fmt.Errorf("error from r.db.Query %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.IsActive); err != nil {
			return []models.User{}, fmt.Errorf("error from rows.Scan %w", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return []models.User{}, fmt.Errorf("error from rows.Err() %w", err)
	}
	return users, nil
}

func (r *repo) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	const query = `
		SELECT id, name, age, is_active 
			FROM users 
		WHERE id = $1;`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.IsActive)
	if err != nil {
		return models.User{}, fmt.Errorf("error from r.db.QueryRow %w", err)
	}
	return user, nil
}

func (r *repo) Delete(ctx context.Context, id int) error {

	const query = `
		UPDATE users 
			SET inActive = true 
		WHERE id = $1;
`
	rows, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("sql.Result.RowsAffected(): %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repo) GetUser(ctx context.Context, id int, req models.User) (models.User, error) {
	var user models.User
	const query = `
		SELECT u.id, u.name, u.age, u.is_active, a.password
			FROM users u
			JOIN auth a ON u.id = a.user_id
		WHERE u.id = $1 AND u.is_active = $2`

	err := r.db.QueryRow(query, id, true).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.IsActive,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errs.ErrNotFound
		}
		if errors.Is(err, errs.ErrIsActiveFalse) {
			return models.User{}, errs.ErrIsActiveFalse
		}
		return models.User{}, fmt.Errorf("query error %w", err)
	}

	return user, nil
}

func (r *repo) UpdatePassword(ctx context.Context, pass models.User) error {
	const query = `
		UPDATE auth 
			SET password = $2
		WHERE user_id = $1`

	rows, err := r.db.Exec(query, pass.ID, pass.NewPassword)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("error from rowsAffected %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *repo) UpdateAge(ctx context.Context, age int, id int) error {

	const query = `
		UPDATE users
			SET age = $2
		WHERE id = $1`

	rows, err := r.db.Exec(query, id, age)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("error from rowsAffected %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}
