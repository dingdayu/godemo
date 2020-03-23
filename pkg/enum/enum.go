package enum

//
const (
	// 中间件将如下字段写入 context
	ContextHeaderAgent = "client_agent"
	ContextClientIp    = "client_ip"
	ContextToken       = "token"
)

// 时间格式
const (
	DateTimeFormat   = "2006-01-02 15:04:05"
	DateTimeHMFormat = "15:04"
	DateFormat       = "2006-01-02"
	TimeFormat       = "15:04:05"
)
