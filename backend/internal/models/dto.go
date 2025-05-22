package models

type MonitorCreateDTO struct {
	Endpoint string `json:"endpoint"`
	Interval int    `json:"interval"`
	Active   bool   `json:"active"`
}

type MonitorUpdateDTO struct {
	Id       string `json:"id"`
	Endpoint string `json:"endpoint"`
	Interval int    `json:"interval"`
	Active   bool   `json:"active"`
}
