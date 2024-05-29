package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Item struct {
	Title string
	URL   string
	Host  string
}

func New(configFuncs ...ConfigFunc) *Client {
	c := &config{}

	for _, configFunc := range configFuncs {
		configFunc(c)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   time.Second * 5,
			Transport: c.roundTripper,
		},
	}
}

type Client struct {
	httpClient *http.Client
}

func (c Client) TopStories(ctx context.Context, maxItems int) ([]Item, error) {
	itemIDs, err := c.topStories(ctx)
	if err != nil {
		return nil, err
	}

	itemIDs = itemIDs[0:maxItems]

	items := []Item{}
	for _, itemID := range itemIDs {
		item, err := c.getItem(ctx, itemID)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (c Client) topStories(ctx context.Context) ([]int, error) {
	var ids []int

	err := c.request(ctx, "/topstories.json", &ids)
	if err != nil {
		return []int{}, err
	}

	return ids, nil
}

func (c Client) getItem(ctx context.Context, id int) (Item, error) {
	item := Item{}
	path := fmt.Sprintf("/item/%d.json", id)

	err := c.request(ctx, path, &item)
	if err != nil {
		return Item{}, err
	}

	parsedURL, err := url.Parse(item.URL)
	if err != nil {
		return Item{}, err
	}

	item.Host = strings.TrimPrefix(parsedURL.Hostname(), "www.")

	return item, nil
}

func (c Client) request(ctx context.Context, path string, target interface{}) error {
	fullURL := fmt.Sprintf("%s%s", "https://hacker-news.firebaseio.com/v0", path)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, http.NoBody)
	if err != nil {
		return err
	}

	r, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)

	err = dec.Decode(target)
	if err != nil {
		return err
	}

	return nil
}
