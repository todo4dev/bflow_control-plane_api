package builder

import "encoding/json"

type Result[TEntity any] struct {
	Offset int64     `json:"offset"`
	Limit  int64     `json:"limit"`
	Total  int64     `json:"total"`
	Items  []TEntity `json:"items"`
}

func NewResultFromRaw[TEntity any](raw *Result[json.RawMessage]) (*Result[TEntity], error) {
	if raw == nil {
		return nil, nil
	}

	items := make([]TEntity, len(raw.Items))
	for i, msg := range raw.Items {
		if err := json.Unmarshal(msg, &items[i]); err != nil {
			return nil, err
		}
	}

	return &Result[TEntity]{
		Offset: raw.Offset,
		Limit:  raw.Limit,
		Total:  raw.Total,
		Items:  items,
	}, nil
}
