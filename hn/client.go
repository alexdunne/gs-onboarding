package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseUrl string
}

type ClientOption func(c *Client)

func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *Client) {
		c.baseUrl = baseUrl
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseUrl: "https://hacker-news.firebaseio.com/v0",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) FetchTopStories() ([]int, error) {
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

func (c *Client) FetchItem(id int) (*Item, error) {
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
