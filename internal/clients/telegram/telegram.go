package telegram

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod       = "getUpdates"
	sendMessageMethod      = "sendMessage"
	sendPhotoMethod        = "sendPhoto"
	editMessageReplyMarkup = "editMessageReplyMarkup"
	deleteMessageMethod    = "deleteMessage"
	answerCallbackMethod   = "answerCallbackQuery"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doQuery(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doQuery(sendMessageMethod, q)
	if err != nil {
		return errors.Wrap(err, "can't send message")
	}

	return nil
}

func (c *Client) SendPhotoWithInlineKeyboard(chatID int, text string, photoURL string, buttons []InlineKeyboardButton) error {
	buttonsMarkup := c.keyboardMarkup(buttons, 2)
	data, err := json.Marshal(buttonsMarkup)
	if err != nil {
		return errors.Wrap(err, "marshaling markup")
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", photoURL)
	q.Add("caption", text)
	q.Add("reply_markup", string(data))

	_, err = c.doQuery(sendPhotoMethod, q)
	if err != nil {
		return errors.Wrap(err, "can't send photo")
	}

	return nil
}

func (c *Client) SendMessageWithInlineKeyboard(chatID int, text string, buttons []InlineKeyboardButton) error {
	buttonsMarkup := c.keyboardMarkup(buttons, 2)
	data, err := json.Marshal(buttonsMarkup)
	if err != nil {
		return errors.Wrap(err, "can't serialise buttons markup")
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("reply_markup", string(data))

	_, err = c.doQuery(sendMessageMethod, q)

	if err != nil {
		return errors.Wrap(err, "can't send buttons markup")
	}

	return nil
}

func (c *Client) EditMessageReplyMarkup(chatID int, messageID int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))
	q.Add("reply_markup", "{}")

	_, err := c.doQuery(editMessageReplyMarkup, q)
	if err != nil {
		return errors.Wrap(err, "can't edit message")
	}

	return nil
}

func (c *Client) DeleteMessage(chatID int, messageID int) (bool, error) {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))

	data, err := c.doQuery(deleteMessageMethod, q)
	if err != nil {
		return false, errors.Wrap(err, "error delete message")
	}

	var success DeleteResponse
	if err := json.Unmarshal(data, &success); err != nil {
		return false, errors.Wrap(err, "error delete message")
	}

	return success.Result, nil
}

func (c *Client) AnswerCallbackQuery(queryID string, text string) error {
	q := url.Values{}
	q.Add("callback_query_id", queryID)
	q.Add("text", text)

	_, err := c.doQuery(answerCallbackMethod, q)
	if err != nil {
		return errors.Wrap(err, "can't answer on callback query")
	}

	return nil
}

func (c *Client) doQuery(method string, query url.Values) ([]byte, error) {
	const errMessage = "can't do request"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	return body, nil
}

func (c *Client) keyboardMarkup(buttons []InlineKeyboardButton, rowLength int) InlineKeyboardMarkup {
	rows := len(buttons) / rowLength
	if len(buttons)%rowLength != 0 {
		rows++
	}

	markup := make([][]InlineKeyboardButton, rows)
	for i := range markup {
		markup[i] = make([]InlineKeyboardButton, 0, rowLength)
	}

	for i, button := range buttons {
		ind := i / rowLength
		markup[ind] = append(markup[ind], button)
	}

	return InlineKeyboardMarkup{InlineKeyboard: markup}
}
