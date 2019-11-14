package health

import (
	"github.com/prometheus/client_golang/prometheus"
)

// https://blog.csdn.net/u014029783/article/details/80001251
type ClusterManager struct {
	Zone             string
	SenderQueueCount *prometheus.Desc
	TimerCount       *prometheus.Desc
	SmsQueueCount    *prometheus.Desc
}

// Simulate prepare the data
func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
	senderQueueCount int64, timerCount int64, smsQueueCount int64,
) {
	//rds := redis.Redis()
	//// 推送队列长度
	//push := stream.NewStreamQueue(rds)
	//senderQueueCount, _ = push.Len()
	//
	//t := timer.NewTimer(rds)
	//timerCount, _ = t.Len()
	//
	//sms := sms2.NewSmsQueue(rds)
	//smsQueueCount, _ = sms.Len()
	return
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.SenderQueueCount
	ch <- c.TimerCount
}

func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
	//senderQueueCount, timerCount, smsQueueCount := c.ReallyExpensiveAssessmentOfTheSystemState()
	//ch <- prometheus.MustNewConstMetric(
	//	c.SenderQueueCount,
	//	prometheus.CounterValue,
	//	float64(senderQueueCount),
	//)
	//ch <- prometheus.MustNewConstMetric(
	//	c.TimerCount,
	//	prometheus.CounterValue,
	//	float64(timerCount),
	//)
	//ch <- prometheus.MustNewConstMetric(
	//	c.SmsQueueCount,
	//	prometheus.CounterValue,
	//	float64(smsQueueCount),
	//)
}

// NewClusterManager creates the two Descs OOMCountDesc and RAMUsageDesc. Note
// that the zone is set as a ConstLabel. (It's different in each instance of the
// ClusterManager, but constant over the lifetime of an instance.) Then there is
// a variable label "host", since we want to partition the collected metrics by
// host. Since all Descs created in this way are consistent across instances,
// with a guaranteed distinction by the "zone" label, we can register different
// ClusterManager instances with the same registry.
func NewClusterManager(zone string) *ClusterManager {
	return &ClusterManager{
		Zone: zone,
		SenderQueueCount: prometheus.NewDesc(
			"message_queue_count",
			"Message Center Timer Count.",
			[]string{},
			prometheus.Labels{"zone": zone},
		),
		TimerCount: prometheus.NewDesc(
			"timer_count",
			"Message Center Timer Count.",
			[]string{},
			prometheus.Labels{"zone": zone},
		),
		SmsQueueCount: prometheus.NewDesc(
			"sms_queue_count",
			"Sms Send Count.",
			[]string{},
			prometheus.Labels{"zone": zone},
		),
	}
}
