package workflows

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// URL - Default Cased API URL
const URL string = "https://api.cased.com"

// Client -
type Client struct {
	URL        string
	HTTPClient *http.Client
	Token      string
}

// NewClient -
func NewClient(url, token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("missing workflow token")
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		URL:        URL,
		Token:      token,
	}

	if url != "" {
		c.URL = url
	}

	return &c, nil
}

var ErrNotFound = errors.New("not found")

type APIError struct {
	Resource string `json:"resource"`
	Path     string `json:"path"`
	Code     string `json:"code"`
}

type APIErrorResponse struct {
	Error   string     `json:"error,omitempty"`
	Errors  []APIError `json:"errors,omitempty"`
	Message string     `json:"message,omitempty"`
}

func (er APIErrorResponse) GetError() error {
	switch er.Error {
	case "not_found":
		return ErrNotFound
	default:
		data, err := json.Marshal(er)
		if err != nil {
			return errors.New("unknown error")
		}

		return errors.New(string(data))
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		errResp := &APIErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
		}

		return nil, errResp.GetError()
	}

	return body, err
}
