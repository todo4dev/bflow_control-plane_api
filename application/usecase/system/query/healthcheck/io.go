package healthcheck

type Query struct {
}

type Result struct {
	Database *string `json:"database"`
	Cache    *string `json:"cache"`
	Storage  *string `json:"storage"`
	Stream   *string `json:"stream"`
}
