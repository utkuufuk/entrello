package config

const (
	PeriodTypeDefault = "default"
	PeriodTypeDay     = "day"
	PeriodTypeHour    = "hour"
	PeriodTypeMinute  = "minute"
)

type Period struct {
	Type     string `json:"type"`
	Interval int    `json:"interval"`
}

type Source struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Strict   bool   `json:"strict"`
	Label    string `json:"label_id"`
	List     string `json:"list_id"`
	Period   Period `json:"period"`
}
