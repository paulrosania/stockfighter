package api

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

type VenuesService struct {
	client *Client
}

type Execution struct {
	Ok       bool      `json:"ok"`
	Account  string    `json:"account"`
	Venue    string    `json:"venue"`
	Symbol   string    `json:"symbol"`
	Order    Order     `json:"order"`
	Price    int       `json:"price"`
	Filled   int       `json:"filled"`
	FilledAt time.Time `json:"filledAt"`

	StandingID int `json:"standingId"`
	IncomingID int `json:"incomingId"`

	StandingComplete bool `json:"standingComplete"`
	IncomingComplete bool `json:"incomingComplete"`
}

type QuoteHandler func(*Quote)

func (s *VenuesService) Quotes(account, venue string, callback QuoteHandler) error {
	path := fmt.Sprintf("%s/venues/%s/tickertape", account, venue)
	conn, err := s.client.DialWebsocket(path)
	if err != nil {
		return err
	}

	for {
		resp := new(QuoteResponse)
		err = websocket.JSON.Receive(conn, resp)
		if err != nil {
			return err
		} else {
			resp.Quote.Ok = resp.Ok
			callback(resp.Quote)
		}
	}
}

type ExecutionHandler func(*Execution)

func (s *VenuesService) Executions(account, venue string, callback ExecutionHandler) error {
	path := fmt.Sprintf("%s/venues/%s/executions", account, venue)
	conn, err := s.client.DialWebsocket(path)
	if err != nil {
		return err
	}

	for {
		execution := new(Execution)
		err = websocket.JSON.Receive(conn, execution)
		if err != nil {
			return err
		} else {
			callback(execution)
		}
	}
}
