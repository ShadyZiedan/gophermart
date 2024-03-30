package repositories

import (
	"context"

	"github.com/ShadyZiedan/gophermart/internal/models"
)

type UsersRepository struct {
	conn pgConn
}

func NewUsersRepository(conn pgConn) *UsersRepository {
	return &UsersRepository{conn: conn}
}

func (ur *UsersRepository) SaveUser(ctx context.Context, user *models.User) error {
	sql := "INSERT INTO users (username, password) values ($1, $2) returning id"
	return ur.conn.QueryRow(ctx,
		sql,
		user.Username,
		user.Password,
	).Scan(&user.Id)
}

func (ur *UsersRepository) FindUserById(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	sql := `SELECT users.id, users.username, users.password FROM users WHERE users.id = $1`
	err := ur.conn.QueryRow(ctx, sql, id).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UsersRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	sql := `SELECT users.id, users.username, users.password FROM users WHERE users.username = $1`
	err := ur.conn.QueryRow(ctx, sql, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UsersRepository) IsUserExist(ctx context.Context, username string) (bool, error) {
	sql := `SELECT exists(SELECT id FROM users WHERE username = $1)`
	var exists bool
	err := ur.conn.QueryRow(ctx, sql, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
