package models

import "fmt"

type MonitorStatus string

const (
	StatusUp   MonitorStatus = "Up"
	StatusDown MonitorStatus = "Down"
)

func (s MonitorStatus) String() string {
	return string(s)
}

func (s MonitorStatus) IsValid() bool {
	switch s {
	case StatusUp, StatusDown:
		return true
	}
	return false
}

type MonitorAction string

const (
	MonitorCreate MonitorAction = "Create"
	MonitorDelete MonitorAction = "Delete"
	MonitorUpdate MonitorAction = "Update"
)

func (s MonitorAction) String() string {
	return string(s)
}

func (s MonitorAction) IsValid() bool {
	switch s {
	case MonitorCreate, MonitorDelete, MonitorUpdate:
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

type MonitorEvent struct {
	Action    MonitorAction `json:"action"`
	MonitorId string        `json:"monitorId"`
}

type MonitorCache struct {
	Endpoint string        `json:"endpoint"`
	Interval int           `json:"interval"`
	Status   MonitorStatus `json:"status"`
}

func (m MonitorCache) String() string {
	return fmt.Sprintf("🖥️ MonitorCache: { Endpoint: %s, Status: %s, Interval: %v }", m.Endpoint, m.Status.String(), m.Interval)
}

type MonitorResult struct {
	Id      string        `json:"id"`
	Date    string        `json:"date"`
	Status  MonitorStatus `json:"status"`
	Latency int64         `json:"latency"`
	Code    int           `json:"code"`
	Error   string        `json:"error"`
}
