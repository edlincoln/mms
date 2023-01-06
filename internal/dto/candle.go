package dto

type Candle struct {
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Open      float64 `json:"open"`
	Timestamp int64   `json:"timestamp"`
	Volume    float64 `json:"volume"`
}

type ExtractorResponse struct {
	StatusCode    int      `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Timestamp     int64    `json:"server_unix_timestamp"`
	Candles       []Candle `json:"candles"`
}
