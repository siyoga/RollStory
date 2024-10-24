package telegram

import (
	"encoding/json"

	"github.com/siyoga/rollstory/internal/errors"
)

func (a *adapter) getUpdates(offset int, limit int) (upd []Update, err error) {
	defer func() {
		if err != nil {
			a.logger.AdapterError(err, errors.ErrTelegramGetUpdates)
		}
	}()

	req, err := json.Marshal(UpdateParam{Offset: offset, Limit: limit})
	if err != nil {
		return nil, err
	}

	data, err := a.makeRequest(methodGetUpdates, req)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (a *adapter) Updates() UpdatesChan {
	ch := make(chan Update, a.batchSize)

	go func() {
		for {
			select {
			case <-a.shutdownChan:
				close(ch)
				return
			default:
			}

			updates, err := a.getUpdates(a.offset, a.batchSize)
			if err != nil {
				a.logger.AdapterError(err, errors.ErrTelegramGetUpdates)
				continue
			}

			for _, update := range updates {
				if update.ID >= a.offset {
					a.offset = update.ID + 1
					ch <- update
				}
			}
		}
	}()

	return ch
}
