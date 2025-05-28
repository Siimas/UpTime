package constants

const (
	RedisMonitorKey          = "monitor"
	RedisMonitorsScheduleKey = "monitors_schedule"

	KafkaMonitorResultsTopic  = "monitor_results"
	KafkaMonitorScheduleTopic = "monitor_schedule"

	PingChanSize        = 1000
	PingWorkerCount     = 10
	ScheduleChanSize    = 1000
	ScheduleWorkerCount = 10

	LoggerChanSize    = 1000
	LoggerWorkerCount = 10
)
