package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/AlexCorn999/link-saving-tgbot/internal/config"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(host string, config *config.Config) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(config.TgToken),
		client:   http.Client{},
	}
}

// newBasePath generates the path for the Telegram request.
func newBasePath(token string) string {
	return "bot" + token
}

// Updates gets a list of updates from Telegram.
func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest("getUpdates", q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	return res.Result, nil
}

// SendMessage sends a message to Telegram.
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	if _, err := c.doRequest("sendMessage", q); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

// doRequest generates and sends a request to Telegram.
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return body, nil
}
