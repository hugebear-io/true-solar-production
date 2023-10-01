package constant

import "time"

const (
	LOW_PERFORMANCE_ALARM = "low_performance_alarm"
	SUM_PERFORMANCE_ALARM = "sum_performance_alarm"
)

const (
	PERFORMANCE_ALARM_SNMP_BATCH_SIZE  = 25
	PERFORMANCE_ALARM_SNMP_BATCH_DELAY = 5 * time.Second
)

const (
	PERFORMANCE_ALARM_TYPE_PERFORMANCE_LOW = iota + 1
	PERFORMANCE_ALARM_TYPE_SUM_PERFORMANCE_LOW
)
