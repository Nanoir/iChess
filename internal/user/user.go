package user

//user.go
import (
	"context"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar,omitempty"`
}

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserReq struct {
	Username string `json:"Username"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateAvatarReq struct {
	ID     int64  `json:"id"`
	Avatar string `json:"avatar"`
}

type VerifyTokenRes struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_token"`
}
type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(c context.Context, userID int64) (*User, error)
	UpdateAvatar(ctx context.Context, req *UpdateAvatarReq) error
	GetAvatar(ctx context.Context, userID int64) (string, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
	VerifyToken(c context.Context, token string) (*VerifyTokenRes, bool, error)
	UpdateAvatar(ctx context.Context, req *UpdateAvatarReq) error
	GetAvatar(ctx context.Context, userID int64) (string, error)
}
