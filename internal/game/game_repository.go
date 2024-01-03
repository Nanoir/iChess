package game

import (
	"database/sql"
	"sync"
)

type GameRepository struct {
	mu    sync.Mutex
	games []Game
	db    *sql.DB // 新增数据库连接
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) AddGame(game Game) {
	r.mu.Lock()
	defer r.mu.Unlock()
	game.ID = len(r.games) + 1
	r.games = append(r.games, game)
}

func (r *GameRepository) GetFilteredGames(player string, date string, eco string) []Game {
	r.mu.Lock()
	defer r.mu.Unlock()

	var filteredGames []Game
	for _, game := range r.games {
		if (player == "" || game.White == player || game.Black == player) &&
			(date == "" || game.Date == date) &&
			(eco == "" || game.ECO == eco) {
			filteredGames = append(filteredGames, game)
		}
	}

	return filteredGames
}
