package router

import (
	"server/internal/game"
	"server/internal/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitUserRouter(userHandler *user.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173" || origin == "http://localhost:3001"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)
	r.POST("/verify_token", userHandler.VerifyToken)
	r.POST("/update_avatar", userHandler.UpdateAvatar)
}

func InitGameRouter(gameHandler *game.GameHandler) {
	router := gin.Default()

	router.GET("/getFiltered", gameHandler.GetFilteredGamesHandler)
}

func Start(addr string) error {
	return r.Run(addr)
}
