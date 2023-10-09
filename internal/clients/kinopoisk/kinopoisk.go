package kinopoisk

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
	headerToken = `X-API-KEY`
	getByTitle  = "v1.2/movie/search"
	getByID     = "v1.3/movie"
)

type KpAPI struct {
	host   string
	token  string
	client *http.Client
}

func NewKp(host string, token string) *KpAPI {
	return &KpAPI{
		host:   host,
		token:  token,
		client: http.DefaultClient,
	}
}

func (k *KpAPI) FindMovieByTitle(title string, limit int) ([]MovieByTitle, error) {
	q := url.Values{}
	q.Add("page", "1")
	q.Add("limit", strconv.Itoa(limit))
	q.Add("query", title)

	data, err := k.doTitleQuery(getByTitle, q)
	if err != nil {
		return nil, errors.Wrap(err, "can't find movie")
	}

	var movies MovieResponse
	if err := json.Unmarshal(data, &movies); err != nil {
		return nil, errors.Wrap(err, "can't find movie")
	}

	return movies.Docs, nil
}

func (k *KpAPI) GetMovieByID(movieID int) (MovieByID, error) {
	data, err := k.doIdQuery(getByID, movieID)
	if err != nil {
		return MovieByID{}, errors.Wrap(err, "can't get movie from api")
	}

	var movie MovieByID
	if err := json.Unmarshal(data, &movie); err != nil {
		return MovieByID{}, errors.Wrap(err, "can't unmarshal json")
	}

	return movie, nil
}

func (k *KpAPI) doTitleQuery(method string, query url.Values) ([]byte, error) {
	const errMessage = "can't do request"

	u := url.URL{
		Scheme: "https",
		Host:   k.host,
		Path:   method,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	req.URL.RawQuery = query.Encode()
	req.Header.Set(`accept`, "application/json")
	req.Header.Add(headerToken, k.token)

	resp, err := k.client.Do(req)
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

func (k *KpAPI) doIdQuery(method string, movieID int) ([]byte, error) {
	const errMessage = "can't do request"

	u := url.URL{
		Scheme: "https",
		Host:   k.host,
		Path:   path.Join(method, strconv.Itoa(movieID)),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Add(headerToken, k.token)

	resp, err := k.client.Do(req)
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
