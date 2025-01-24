package telegram

import (
	"encoding/json"
	"go-bot/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(token string) *Client {
	return &Client{
		host:     "https://api.telegram.org",
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "/bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", string(offset))
	q.Add("limit", string(limit))

	// do request
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, e.Wrap(err, "failed to do request")
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap(err, "failed to unmarshal response")
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chat int, text string) error {
	q := url.Values{}
	q.Add("chat_id", string(chat))
	q.Add("text", string(text))

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		e.Wrap(err, "failed to send message")
	}

	return nil
}

func (c *Client) doRequest(method string, q url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(err, "failed to create request")
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(err, "failed to do request")
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(err, "failed to read body")
	}
	return body, nil
}
