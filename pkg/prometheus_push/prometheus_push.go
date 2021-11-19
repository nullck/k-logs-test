package prometheus_push

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PrometheusPusher struct {
	GWUrl      string
	GWPort     int
	MetricName string
}

func (p PrometheusPusher) PushMetric(delayPeriod int64) {
	//ref: https://godoc.org/github.com/prometheus/client_golang/prometheus/push#Pusher.Add
	delayTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "logs_delay_in_milliseconds",
		Help: "Logs delay time defined in ms",
	})
	promoGw := fmt.Sprintf("%v:%v", p.GWUrl, p.GWPort)
	registry := prometheus.NewRegistry()
	registry.MustRegister(delayTime)
	pusher := push.New(promoGw, p.MetricName).Gatherer(registry)
	delayTime.Set(float64(delayPeriod))
	if err := pusher.Add(); err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	}
}
