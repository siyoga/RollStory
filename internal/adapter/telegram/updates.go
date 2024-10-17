package telegram

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/siyoga/rollstory/internal/domain"
)

func (a *adapter) Updates(offset int, limit int) (updates []domain.Update, err error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := a.makeRequest(methodGetUpdates, q)
	if err != nil {
		return nil, err
	}

	var res domain.UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}
