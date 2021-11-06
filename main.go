package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"html/template"
	"net/http"
	"os"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		}, []string{"code", "method", "endpoint"})

	prometheusRegistry = prometheus.NewRegistry()
)

type GetIndexHandler struct {
	HelloClient *HelloClient
}

func GetTraceHeadersFromRequest(req *http.Request) map[string]string {
	headers := make(map[string]string)
	relevantHeaderNames := []string{
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-ot-span-context",
	}

	for _, header := range relevantHeaderNames {
		if value := req.Header.Get(header); value != "" {
			headers[header] = value
		}
	}

	return headers
}

func (h *GetIndexHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	message, _ := h.HelloClient.GetHello("There", GetTraceHeadersFromRequest(req))

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, message)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	//Setup Prometheus
	prometheusRegistry.MustRegister(httpRequestsTotal)

	//Configure HTTP server
	http.Handle("/metrics", promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/",
		promhttp.InstrumentHandlerCounter(
			httpRequestsTotal.MustCurryWith(prometheus.Labels{"endpoint": "index"}),
			&GetIndexHandler{HelloClient: NewHelloClient(getEnv("HELLOAPP_URL", "http://localhost:8080"))},
		),
	)
	http.ListenAndServe(":"+getEnv("PORT", "8080"), nil)
}
