package handlers

type ServerEvent struct {
	EventId   string
	Title     string
	Message   string
	Level     string
	Timestamp string
}

type ServerUsage struct {
	CPU_MAX_USAGE        float64
	RAM_MAX_USAGE        float64
	DISK_MAX_USAGE       float64
	USAGE_INTERVAL_CHECK int
	TIMESTAMP            string
	ERROR                string
}
