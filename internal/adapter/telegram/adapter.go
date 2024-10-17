package telegram

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	def "github.com/siyoga/rollstory/internal/adapter"
	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/logger"
)

var _ def.TelegramAdapter = (*adapter)(nil)

const (
	methodGetUpdates  = "getUpdates"
	methodSendMessage = "sendMessage"
)

type (
	adapter struct {
		host     string
		basePath string

		logger logger.Logger
		client http.Client
	}
)

func NewAdapter(
	cfg config.Bot,

	logger logger.Logger,
) *adapter {
	return &adapter{
		host:     cfg.Host,
		basePath: fmt.Sprintf("bot%s", cfg.Token),

		logger: logger,
		client: http.Client{},
	}
}

func (a *adapter) makeRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		if err != nil {
			err = a.logger.AdapterError(err, nil, method)
		}
	}()

	u := url.URL{
		Scheme: "https",
		Host:   a.host,
		Path:   path.Join(a.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
