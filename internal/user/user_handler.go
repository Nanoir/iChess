package user

// user_handler.go
import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) VerifyToken(c *gin.Context) {
	token, err := getTokenFromHeader(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	ParsedToken, IsValid, err := h.Service.VerifyToken(c.Request.Context(), token)
	if err != nil || !IsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	fmt.Println(ParsedToken)

	c.JSON(http.StatusOK, ParsedToken)
}

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header format")
	}

	return parts[1], nil
}

func (h *Handler) Login(c *gin.Context) {
	var user LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Service.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将 Access Token 和 Refresh Token 发送到客户端
	c.SetCookie("AccessToken", u.AccessToken, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("RefreshToken", u.RefreshToken, 7*24*60*60, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"username":     u.Username,
		"userID":       u.ID,
		"accessToken":  u.AccessToken,
		"refreshToken": u.RefreshToken})
}

func (h *Handler) Logout(c *gin.Context) {
	// 清除 Access Token 的 Cookie
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)

	// 清除 Refresh Token 的 Cookie
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h *Handler) UpdateAvatar(c *gin.Context) {
	var req UpdateAvatarReq
	fmt.Println("fuck")
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(c)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 在实际应用中，可能需要将前端上传的 base64 数据进行解码，这里简化为直接存储
	err := h.Service.UpdateAvatar(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully"})
}

func (h *Handler) GetAvatar(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	avatar, err := h.Service.GetAvatar(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar": avatar})
}
