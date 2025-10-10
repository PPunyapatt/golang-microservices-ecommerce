package monitor

import (
	"log"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	goroutinesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_goroutines_cart",
		Help: "Number of running goroutines",
	})
)

func StartGoRoutineMonitor(prometheusUrl, jobName string, interval time.Duration) {
	go func() {
		for {
			log.Println("number of goroutines:", runtime.NumGoroutine())
			goroutinesGauge.Set(float64(runtime.NumGoroutine()))
			if err := push.New(prometheusUrl, jobName).
				Collector(goroutinesGauge).
				Push(); err != nil {
				log.Println("Could not push to PushGateway:", err.Error())
			}
			time.Sleep(interval)
		}
	}()
}
