package api

import (
	"fmt"
	"time"
)

type OrdersService struct {
	client *Client
}

type OrderParams struct {
	Account   string `json:"account"`
	Price     int    `json:"price,omitempty"`
	Quantity  int    `json:"qty"`
	Direction string `json:"direction"`
	OrderType string `json:"orderType"`
}

type Order struct {
	Ok               bool      `json:"ok"`
	Symbol           string    `json:"symbol"`
	Venue            string    `json:"venue"`
	Direction        string    `json:"direction"`
	OriginalQuantity int       `json:"originalQty"`
	Quantity         int       `json:"qty"`
	Price            int       `json:"price"`
	Type             string    `json:"string"`
	ID               int       `json:"id"`
	Account          string    `json:"account"`
	ReceivedAt       time.Time `json:"ts"`
	Fills            []Fill    `json:"fills"`
	TotalFilled      int       `json:"totalFilled"`
	Open             bool      `json:"open"`
}

type Fill struct {
	Price    int       `json:"price"`
	Quantity int       `json:"qty"`
	FilledAt time.Time `json:"ts"`
}

type OrderList struct {
	Ok     bool     `json:"ok"`
	Venue  string   `json:"venue"`
	Orders []*Order `json:"orders"`
}

func (s *OrdersService) New(venue, symbol string, params *OrderParams) (*Order, error) {
	var r Order
	path := fmt.Sprintf("venues/%s/stocks/%s/orders", venue, symbol)
	err := s.client.Call("POST", path, params, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *OrdersService) Get(venue, symbol string, id int) (*Order, error) {
	var r Order
	path := fmt.Sprintf("venues/%s/stocks/%s/orders/%d", venue, symbol, id)
	err := s.client.Call("GET", path, nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *OrdersService) All(venue, account string) ([]*Order, error) {
	var r OrderList
	path := fmt.Sprintf("venues/%s/accounts/%s/orders", venue, account)
	err := s.client.Call("GET", path, nil, &r)
	if err != nil {
		return nil, err
	}

	return r.Orders, nil
}

func (s *OrdersService) BySymbol(venue, account, symbol string) ([]*Order, error) {
	var r OrderList
	path := fmt.Sprintf("venues/%s/accounts/%s/stocks/%s/orders", venue, account, symbol)
	err := s.client.Call("GET", path, nil, &r)
	if err != nil {
		return nil, err
	}

	return r.Orders, nil
}

func (s *OrdersService) Cancel(venue, symbol string, id int) (*Order, error) {
	var r Order
	path := fmt.Sprintf("venues/%s/stocks/%s/orders/%d", venue, symbol, id)
	err := s.client.Call("DELETE", path, nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
