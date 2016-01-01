package api

import (
	"fmt"
)

type InstancesService struct {
	client *Client
}

type InstanceDetails struct {
	EndOfTheWorldDay int `json:"endOfTheWorldDay"`
	TradingDay       int `json:"tradingDay"`
}

type Instance struct {
	Ok      bool              `json:"ok"`
	Done    bool              `json:"done"`
	ID      int               `json:"id"`
	State   string            `json:"state"`
	Details InstanceDetails   `json:"details,omitempty"`
	Flash   map[string]string `json:"flash,omitempty"`
}

func (s *InstancesService) Get(id int) (*Instance, error) {
	var r Instance
	path := fmt.Sprintf("/gm/instances/%d", id)
	err := s.client.Call("GET", path, nil, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
