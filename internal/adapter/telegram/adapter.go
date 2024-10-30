package telegram

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/logger"
)

const (
	methodGetUpdates = "getUpdates"

	methodSendMessage   = "sendMessage"
	methodDeleteMessage = "deleteMessage"
	methodEditMessage   = "editMessageText"
)

type (
	Adapter interface {
		SendMessage(req Response) (Message, error)
		DeleteMessage(chatId int, messageId int) (err error)
		EditMessage(chatId int, messageId int, text string, markups *domain.InlineMarkup) (err error)
		Updates() UpdatesChan
	}

	adapter struct {
		host      string
		basePath  string
		batchSize int
		offset    int

		logger logger.Logger

		client       http.Client
		shutdownChan chan struct{}
	}
)

func NewAdapter(
	cfg config.Bot,

	logger logger.Logger,
) Adapter {
	return &adapter{
		host:      cfg.Host,
		basePath:  fmt.Sprintf("bot%s", cfg.Token),
		offset:    cfg.Offset,
		batchSize: cfg.BatchSize,

		logger:       logger,
		client:       http.Client{},
		shutdownChan: make(chan struct{}),
	}
}

func (a *adapter) makeRequest(method string, reqBody []byte) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   a.host,
		Path:   path.Join(a.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

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
