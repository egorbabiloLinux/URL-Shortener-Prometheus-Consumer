package metrics 

import (
	"github.com/prometheus/client_golang/prometheus"
)


var (
	LoginCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:"app_logins_total",
			Help:"Total number of logins events"
		},
		[]string{"source"},
	)
)

func Register(){
	prometheus.MustRegister(LoginCounter)
}
