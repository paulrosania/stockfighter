package api

import (
	"bytes"
	"fmt"
	"io"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"golang.org/x/net/websocket"
)

const DefaultBaseURL = "https://api.stockfighter.io/ob/api/"
const DefaultWebsocketBaseURL = "wss://api.stockfighter.io/ob/api/ws/"
const DefaultWebsocketOrigin = "http://localhost"
const DefaultUserAgent = "stockfighter-go/0.1"

type Client struct {
	client *http.Client

	BaseURL          *url.URL
	WebsocketBaseURL *url.URL
	WebsocketOrigin  *url.URL
	UserAgent        string

	APIKey string

	Games     GamesService
	Instances InstancesService
	Orders    OrdersService
	Stocks    StocksService
	Venues    VenuesService
}

func NewClient(apiKey string) *Client {
	httpu, err := url.Parse(DefaultBaseURL)
	if err != nil {
		panic(err)
	}

	wssu, err := url.Parse(DefaultWebsocketBaseURL)
	if err != nil {
		panic(err)
	}

	wso, err := url.Parse(DefaultWebsocketOrigin)
	if err != nil {
		panic(err)
	}

	client := &Client{
		client:           http.DefaultClient,
		BaseURL:          httpu,
		WebsocketBaseURL: wssu,
		WebsocketOrigin:  wso,
		UserAgent:        DefaultUserAgent,
		APIKey:           apiKey,
	}

	client.Games = GamesService{client: client}
	client.Instances = InstancesService{client: client}
	client.Orders = OrdersService{client: client}
	client.Stocks = StocksService{client: client}
	client.Venues = VenuesService{client: client}

	return client
}

func (c *Client) DialWebsocket(urlStr string) (*websocket.Conn, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.WebsocketBaseURL.ResolveReference(rel)
	conf := &websocket.Config{
		Location: u,
		Origin:   c.WebsocketOrigin,
		Version:  websocket.ProtocolVersionHybi13,
	}
	return websocket.DialConfig(conf)
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var contentType string
	var buf io.ReadWriter
	if body != nil {
		if method == "GET" {
			v, err := query.Values(body)
			if err != nil {
				return nil, err
			}
			u.RawQuery = v.Encode()
		} else {
			buf = new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
			contentType = "application/json"
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Starfighter-Authorization", c.APIKey)
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}
	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) Call(method, urlStr string, body interface{}, v interface{}) error {
	req, err := c.NewRequest(method, urlStr, body)
	if err != nil {
		return err
	}

	err = c.Do(req, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("request failed with empty body and status: %s", resp.Status)
	}

	e := new(ErrorResponse)
	err = json.Unmarshal(body, e)
	if err != nil {
		return fmt.Errorf("request failed with non-JSON body and status: %s", resp.Status)
	}

	return e
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if v != nil {
		err = json.Unmarshal(body, v)
	}

	return err
}
