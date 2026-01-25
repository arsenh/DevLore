package repository

import (
	"context"
	"errors"

	"github.com/arsenh/DevLore/internal/database"
	"github.com/arsenh/DevLore/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	CommandUserFindByID    = "SELECT id, full_name, email, password_hash, created_at FROM users WHERE id = $1;"
	CommandUserFindByEmail = "SELECT id, full_name, email, password_hash, created_at FROM users WHERE email = $1;"
	CommandUserCreate      = "INSERT INTO users (full_name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at;"
	CommandUserList        = "SELECT * FROM users;"
	CommandUserDelete      = "DELETE FROM users WHERE id = $1;"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresUserRepository struct {
	db *database.DbContext
}

func NewPostgresUserRepository(db *database.DbContext) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (d *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {

	var user model.User

	err := d.db.QueryRow(ctx, CommandUserFindByID, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {

	var user model.User

	err := d.db.QueryRow(ctx, CommandUserFindByEmail, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *PostgresUserRepository) Create(ctx context.Context, u *model.User) error {
	return d.db.QueryRow(ctx, CommandUserCreate,
		u.FullName,
		u.Email,
		u.Password,
	).Scan(&u.ID, &u.CreatedAt)
}

func (d *PostgresUserRepository) List(ctx context.Context) ([]model.User, error) {
	rows, err := d.db.QueryRows(ctx, CommandUserList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}

	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (d *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := d.db.Exec(ctx, CommandUserDelete, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
