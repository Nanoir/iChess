package user

//user_repository.go
import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int64
	query := "INSERT INTO users(username, password, email) VALUES (?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Email)
	if err != nil {
		return &User{}, err
	}

	lastInsertId, err = result.LastInsertId()
	if err != nil {
		return &User{}, err
	}

	user.ID = lastInsertId
	return user, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	u := User{}
	query := "SELECT id, email, username, password FROM users WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return &User{}, err
	}

	return &u, nil
}

func (r *repository) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	u := User{}
	query := "SELECT id, email, username, password, avatar FROM users WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return &User{}, err
	}

	return &u, nil
}

func (r *repository) UpdateAvatar(ctx context.Context, req *UpdateAvatarReq) error {
	query := "UPDATE users SET avatar = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, req.Avatar, req.ID)
	return err
}

func (r *repository) GetAvatar(ctx context.Context, userID int64) (string, error) {
	var avatar string
	query := "SELECT avatar FROM users WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("avatar not found")
		}
		return "", err
	}
	return avatar, nil
}
