package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"html/template"
	"net/http"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		}, []string{"code", "method", "endpoint"})

	prometheusRegistry = prometheus.NewRegistry()
)

func GetIndex(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, "Hello World!")
}

func main() {
	//Setup Prometheus
	prometheusRegistry.MustRegister(httpRequestsTotal)

	//Configure HTTP server
	http.Handle("/metrics", promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", promhttp.InstrumentHandlerCounter(httpRequestsTotal.MustCurryWith(prometheus.Labels{"endpoint": "index"}), http.HandlerFunc(GetIndex)))
	http.ListenAndServe(":8080", nil)
}
