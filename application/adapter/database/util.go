package database

import "encoding/json"

func TypedFromJsonWithErr[TEntity any](msg *json.RawMessage, err error) (*TEntity, error) {
	if err != nil {
		return nil, err
	}

	if msg == nil {
		return nil, nil
	}

	var out TEntity
	if err := json.Unmarshal(*msg, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
