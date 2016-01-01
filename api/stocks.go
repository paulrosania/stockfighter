package api

import (
	"fmt"
	"time"
)

type StocksService struct {
	client *Client
}

type Stock struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type StockList struct {
	Ok      bool    `json:"ok"`
	Symbols []Stock `json:"symbols"`
}

type Quote struct {
	Ok bool `json:"ok"`

	Symbol string `json:"symbol"`
	Venue  string `json:"venue"`

	Bid      int `json:"bid"`
	BidSize  int `json:"bidSize"`
	BidDepth int `json:"bidDepth"`

	Ask      int `json:"ask"`
	AskSize  int `json:"askSize"`
	AskDepth int `json:"askDepth"`

	Last        int       `json:"last"`
	LastSize    int       `json:"lastSize"`
	LastTradeAt time.Time `json:"lastTrade"`

	UpdatedAt time.Time `json:"quoteTime"`
}

type QuoteResponse struct {
	Ok    bool   `json:"ok"`
	Quote *Quote `json:"quote"`
}

type OrderbookEntry struct {
	Price    int  `json:"price"`
	Quantity int  `json:"qty"`
	IsBuy    bool `json:"isBuy"`
}

type Orderbook struct {
	Ok bool `json:"ok"`

	Symbol string `json:"symbol"`
	Venue  string `json:"venue"`

	Bids []OrderbookEntry `json:"bids"`
	Asks []OrderbookEntry `json:"asks"`

	UpdatedAt time.Time `json:"ts"`
}

func (s *StocksService) Quote(venue, symbol string) (*Quote, error) {
	var r Quote
	path := fmt.Sprintf("venues/%s/stocks/%s/quote", venue, symbol)
	err := s.client.Call("GET", path, nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *StocksService) Orderbook(venue, symbol string) (*Orderbook, error) {
	var r Orderbook
	path := fmt.Sprintf("venues/%s/stocks/%s", venue, symbol)
	err := s.client.Call("GET", path, nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *StocksService) List(venue string) ([]Stock, error) {
	var r StockList
	path := fmt.Sprintf("venues/%s/stocks", venue)
	err := s.client.Call("GET", path, nil, &r)
	if err != nil {
		return nil, err
	}
	return r.Symbols, nil
}
