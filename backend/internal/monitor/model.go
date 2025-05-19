package monitor

type MonitorStatus string

const (
	StatusUp   MonitorStatus = "Up"
	StatusDown MonitorStatus = "Down"
)

func (s MonitorStatus) string() string {
	return string(s)
}

func (s MonitorStatus) IsValid() bool {
	switch s {
	case StatusUp, StatusDown:
		return true
	}
	return false
}

type Monitor struct {
	Id       string `json:"id"`
	Endpoint string `json:"endpoint"`
	Interval int    `json:"interval"`
	Active   bool   `json:"active"`
}

type MonitorCache struct {
	Endpoint string        `json:"endpoint"`
	Interval int           `json:"interval"`
	Status   MonitorStatus `json:"status"`
}

type MonitorResult struct {
	Id      string        `json:"id"`
	Date    string        `json:"date"`
	Status  MonitorStatus `json:"status"`
	Latency int64         `json:"latency"`
	Code    int           `json:"code"`
	Error   string        `json:"error"`
}
