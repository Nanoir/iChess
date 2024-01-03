package main

import (
	_ "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"server/db"
	"server/internal/game"
	"server/internal/user"
	"server/middleware"
	"server/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}
	RedisConn := middleware.InitRedisClient()

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep, RedisConn)
	userHandler := user.NewHandler(userSvc)

	gameRepo := game.NewGameRepository(dbConn.GetDB())
	gameSvc := game.NewGameService(gameRepo)
	gameHandler := game.NewGameHandler(gameSvc)
	router.InitGameRouter(gameHandler)

	router.InitUserRouter(userHandler)
	err = router.Start("0.0.0.0:8000")
	if err != nil {
		return
	}
}
