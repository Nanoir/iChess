package game

type GameService struct {
	repo *GameRepository
}

func NewGameService(repo *GameRepository) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) AddGame(game Game) {
	s.repo.AddGame(game)
}

func (s *GameService) GetFilteredGames(player string, date string, eco string) []Game {
	return s.repo.GetFilteredGames(player, date, eco)
}
