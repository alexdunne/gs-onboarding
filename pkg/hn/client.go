package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	FetchTopStories() ([]int, error)
	FetchItem(id int) (*Item, error)
}

type client struct {
	baseUrl string
}

type ClientOption func(c *client)

func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *client) {
		c.baseUrl = baseUrl
	}
}

func New(opts ...ClientOption) *client {
	c := &client{
		baseUrl: "https://hacker-news.firebaseio.com/v0",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

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
