package unsplash

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type APIResponse struct {
	Total      int             `json:"total"`
	TotalPages int             `json:"total_pages"`
	Results    []PictureResult `json:"results"`
}

type PictureResult struct {
	ID      string            `json:"id"`
	Width   int               `json:"width"`
	Height  int               `json:"height"`
	URLs    map[string]string `json:"urls"`
	Resized string
}

type APIClient struct {
	// note how the lowercase accessToken is private
	accessToken string
}

func NewAPIClient(token string) APIClient {
	return APIClient{token}
}

func (c *APIClient) Search(query string) (*APIResponse, error) {
	url := fmt.Sprintf("https://api.unsplash.com/search/photos?page=1&query=%s&client_id=%s", query, c.accessToken)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from Unsplash API")
	}
	defer resp.Body.Close()

	response := APIResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON")
	}

	return &response, nil
}

func LoadRockets() (*APIResponse, error) {
	query := "spacex"
	client := NewAPIClient("c1f9504a548ec5ea75acf3a3919ceab1ab04d09b732a839f04ca0be74f6227a0")
	response, err := client.Search(query)
	return response, err
}
