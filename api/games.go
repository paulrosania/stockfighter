package api

import (
	"fmt"
)

type GamesService struct {
	client *Client
}

type Game struct {
	Ok         bool `json:"ok"`
	InstanceID int  `json:"instanceId"`

	Account  string         `json:"account"`
	Tickers  []string       `json:"tickers"`
	Venues   []string       `json:"venues"`
	Balances map[string]int `json:"balances"`

	Instructions map[string]string `json:"instructions"`

	SecondsPerTradingDay int `json:"secondsPerTradingDay"`
}

func (s *GamesService) New(game string) (*Game, error) {
	var r Game
	path := fmt.Sprintf("/gm/levels/%s", game)
	err := s.client.Call("POST", path, nil, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
