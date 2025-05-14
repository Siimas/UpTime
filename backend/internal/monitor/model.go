package monitor

type MonitorStatus string

const (
	StatusOnline  MonitorStatus = "Online"
	StatusOffline MonitorStatus = "Offline"
	StatusPaused  MonitorStatus = "Paused"
)

func (s MonitorStatus) IsValid() bool {
	switch s {
	case StatusOnline, StatusOffline, StatusPaused:
		return true
	}
	return false
}

type Monitor struct {
	Endpoint string `json:"endpoint"`
	Interval int    `json:"interval"`
	Status   MonitorStatus `json:"status"`
}