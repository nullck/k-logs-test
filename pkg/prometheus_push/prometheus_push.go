package prometheus_push

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Prometheus struct {
	GWUrl      string
	GWPort     int
	MetricName string
}

func (p Prometheus) PushMetric(delayPeriod int64) {
	delayTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "logs_delay_in_milliseconds",
		Help: "Logs delay time defined in ms",
	})
	promoGw := fmt.Sprintf("%v:%v", p.GWUrl, p.GWPort)
	if err := push.New(promoGw, p.MetricName).
		Collector(delayTime).
		Grouping("job", "logs_delay").
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
}
