package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/fnuritdinov/firstService/internal/models"
)

type IUserRepo interface {
	Login(ctx context.Context, auth models.Auth) error
	GetAll(ctx context.Context) ([]models.User, error)
	Get(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, id int, updatedUser string) error
	Delete(ctx context.Context, id int) error
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
	const query = `INSERT INTO auth (login, password, user_id) values ($1, $2, $3);`

	rows, err := r.db.Exec(query, auth.Login, auth.Password, auth.UserID)
	if err != nil {
		fmt.Errorf("error from r.db.Exec %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		fmt.Errorf("sql.Result.RowsAffected(): %w", err)
	}

	fmt.Println(rowsAffected)
	return nil

}

func (r *repo) GetAll(ctx context.Context) ([]models.User, error) {

	const query = `SELECT name, age, password, is_active FROM users WHERE id = ANY($1);`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.User{}, fmt.Errorf("error from r.db.Query %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.IsActive); err != nil {
			return []models.User{}, fmt.Errorf("error from rows.Scan %w", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("error from rows.Err()")
	}
	return users, nil
}

func (r *repo) Get(ctx context.Context) ([]models.User, error) {

	const query = `SELECT id, name, age, is_active FROM users WHERE inActive = false;`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.User{}, fmt.Errorf("error from r.db.Query %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.IsActive); err != nil {
			return []models.User{}, fmt.Errorf("error from rows.Scan %w", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("error from rows.Err()")
	}
	return users, nil
}

func (r *repo) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	const query = `SELECT id, name, age, is_active FROM users WHERE id = $1;`
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

func (r *repo) Create(ctx context.Context, user models.User) error {

	const query = `INSERT INTO users (name, age, is_active) VALUES ($1, $2, $3);`

	rows, err := r.db.Exec(query, user.Name, user.Age, user.IsActive)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("sql.Result.RowsAffected(): %w", err)
	}

	fmt.Println(rowsAffected)
	return nil
}

func (r *repo) Update(ctx context.Context, id int, updatedUser string) error {

	const query = `UPDATE users SET name = $1 WHERE id = $2;`

	rows, err := r.db.Exec(query, updatedUser, id)
	if err != nil {
		return fmt.Errorf("error from r.db.Exec %w", rows)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("error from rowsAffected %w", err)
	}
	fmt.Println(rowsAffected)
	return nil
}

func (r *repo) Delete(ctx context.Context, id int) error {

	const query = `UPDATE users SET inActive = true WHERE id = $1;`

	rows, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Errorf("error from r.db.Exec %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("sql.Result.RowsAffected(): %w", err)
	}

	fmt.Println(rowsAffected)
	return nil
}
