package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is a interface to expose methods to interact with the hacker news api
type Client interface {
	FetchTopStories() ([]int, error)
	FetchItem(id int) (*Item, error)
}

type client struct {
	baseUrl string
}

// ClientOption is an interface for a functional option
type ClientOption func(c *client)

// WithBaseUrl is a functional option to configure the client's base URL
func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *client) {
		c.baseUrl = baseUrl
	}
}

// New creates a client
func New(opts ...ClientOption) *client {
	c := &client{
		baseUrl: "https://hacker-news.firebaseio.com/v0",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// FetchTopStories fetches the ids of the current top hacker news stories
func (c *client) FetchTopStories() ([]int, error) {
	resp, err := http.Get(c.baseUrl + "/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res []int
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

// FetchItem fetches item information for a given id from the hacker news api
func (c *client) FetchItem(id int) (*Item, error) {
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", c.baseUrl, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res Item
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
