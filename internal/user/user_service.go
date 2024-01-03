package user

//user_service.go
import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	accessSecretKey  = "access_secret"
	refreshSecretKey = "refresh_secret"
)

type service struct {
	Repository  Repository
	RedisClient *redis.Client
	timeout     time.Duration
}

func NewService(repository Repository, redisClient *redis.Client) Service {
	return &service{
		repository,
		redisClient,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

type MyJWTClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return &LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &LoginUserRes{}, err
	}

	accessTokenClaims := &MyJWTClaims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenSS, err := accessToken.SignedString([]byte(accessSecretKey))
	if err != nil {
		return &LoginUserRes{}, err
	}

	// 生成 Refresh Token
	refreshTokenClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh Token 有效期为7天
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenSS, err := refreshToken.SignedString([]byte(refreshSecretKey))
	if err != nil {
		return &LoginUserRes{}, err
	}

	// 返回 Access Token 和 Refresh Token
	return &LoginUserRes{
		ID:           strconv.Itoa(int(u.ID)),
		Username:     u.Username,
		AccessToken:  accessTokenSS,
		RefreshToken: refreshTokenSS,
	}, nil
}
func (s *service) VerifyToken(c context.Context, token string) (*VerifyTokenRes, bool, error) {
	claims := &MyJWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecretKey), nil
	})

	if err != nil {
		return nil, false, err
	}

	if !parsedToken.Valid {
		return &VerifyTokenRes{}, false, nil
	}

	// 检查是否需要刷新令牌
	if time.Until(parsedToken.Claims.(*MyJWTClaims).ExpiresAt.Time) < 30*time.Minute {
		// 如果距离过期时间小于30分钟，可以进行刷新
		refreshToken := jwt.New(jwt.SigningMethodHS256)
		refreshToken.Claims = jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * 3600).Unix(), // Refresh Token 有效期为7天
		}

		refreshSS, err := refreshToken.SignedString([]byte(refreshSecretKey))
		if err != nil {
			return nil, false, err
		}

		// 返回刷新令牌，客户端可以使用刷新令牌请求新的 Access Token
		return &VerifyTokenRes{
			ID:           claims.ID,
			Username:     claims.Username,
			RefreshToken: refreshSS,
		}, true, nil
	}

	// token 验证通过，返回解析出的用户信息
	return &VerifyTokenRes{
		ID:       claims.ID,
		Username: claims.Username,
	}, true, nil
}

func (s *service) UpdateAvatar(ctx context.Context, req *UpdateAvatarReq) error {
	err := s.Repository.UpdateAvatar(ctx, req)
	return err
}

func (s *service) GetAvatar(ctx context.Context, userID int64) (string, error) {
	avatar, err := s.Repository.GetAvatar(ctx, userID)
	if err != nil {
		return "", err
	}
	return avatar, nil
}

func (s *service) setUserOnline(userID int64) error {
	key := fmt.Sprintf("user:%d:status", userID)
	return s.RedisClient.Set(context.Background(), key, "ServerOn", 0).Err()
}

func (s *service) setUserOffline(userID int64) error {
	key := fmt.Sprintf("user:%d:status", userID)
	return s.RedisClient.Del(context.Background(), key).Err()
}

func (s *service) isUserOnline(userID int64) (bool, error) {
	key := fmt.Sprintf("user:%d:status", userID)
	result, err := s.RedisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil // 用户不在线
	} else if err != nil {
		return false, err // 发生错误
	}

	return result == "ServerOn", nil // 用户在线
}
