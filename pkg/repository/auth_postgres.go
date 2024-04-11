package repository

import (
	"fmt"
	"github.com/igostfost/avito_backend_trainee/pkg/types"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user types.User) (int, error) {
	var userId int
	query := fmt.Sprintf("INSERT INTO \"%s\" (username, password_hash) VALUES ($1, $2) RETURNING user_id", usersTable)
	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

func (r *AuthPostgres) CreateAdmin(user types.User) (int, bool, error) {
	var userId int
	var isAdmin bool
	query := fmt.Sprintf("INSERT INTO \"%s\" (username, password_hash, is_admin) VALUES ($1,$2,$3) RETURNING user_id,is_admin", usersTable)
	row := r.db.QueryRow(query, user.Username, user.Password, true)
	if err := row.Scan(&userId, &isAdmin); err != nil {
		return 0, false, err
	}
	return userId, isAdmin, nil
}

func (r *AuthPostgres) GetUser(username, password string) (types.User, error) {
	var user types.User
	query := fmt.Sprintf("SELECT user_id, is_admin FROM \"%s\" WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
