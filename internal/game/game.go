package game

type Game struct {
	ID     int    `json:"id"`
	White  string `json:"white"`
	EloW   string `json:"elo_w"`
	Black  string `json:"black"`
	EloB   string `json:"elo_b"`
	Result string `json:"result"`
	ECO    string `json:"eco"`
	Date   string `json:"date"`
}
