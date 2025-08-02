package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)


var (
	AuthCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:"app_auth_total",
			Help:"Total number of auth events",
		},
		[]string{"source"},
	)
)

func Register(){
	prometheus.MustRegister(AuthCounter)
}
