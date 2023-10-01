package constant

const LOG_DIRECTORY = "logs"

const (
	DAILY_PRODUCTION_LOG_NAME      = "daily_production.log"
	MONTHLY_PRODUCTION_LOG_NAME    = "monthly_production.log"
	SOLARMAN_COLLECTOR_LOG_NAME    = "solarman_collector.log"
	SOLARMAN_ALARM_LOG_NAME        = "solarman_alarm.log"
	LOW_PERFORMANCE_ALARM_LOG_NAME = "low_performance_alarm.log"
)

func GetLogName(logName string) string {
	return LOG_DIRECTORY + "/" + logName
}
