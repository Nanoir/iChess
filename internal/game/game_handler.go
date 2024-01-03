package game

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GameHandler struct {
	service *GameService
}

func NewGameHandler(service *GameService) *GameHandler {
	return &GameHandler{service: service}
}

func (h *GameHandler) GetFilteredGamesHandler(c *gin.Context) {
	player := c.Query("player")
	date := c.Query("date")
	eco := c.Query("eco")

	filteredGames := h.service.GetFilteredGames(player, date, eco)

	c.JSON(http.StatusOK, filteredGames)
}
