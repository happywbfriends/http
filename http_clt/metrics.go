package http_clt

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	metricsNamespace           = "sapi"
	metricsSubsystemHttpClient = "http_clt"
	metricsLabelClient         = "client"
	metricsLabelRequest        = "request"
)

type IHttpClientMetrics interface {
	RequestDuration(duration time.Duration)
	IncSuccess()
	IncError()
	IncBadStatus(status int)
}

type NoHttpClientMetrics struct{}

func (m *NoHttpClientMetrics) RequestDuration(time.Duration) {}
func (m *NoHttpClientMetrics) IncSuccess()                   {}
func (m *NoHttpClientMetrics) IncError()                     {}
func (m *NoHttpClientMetrics) IncBadStatus(int)              {}

/*
	Метрики следует считать в разрезе конкретного клиента и конкретного запроса.
	Например, в нашем сервисе есть 2 http клиента: один ходит в ABAC, другой - в SuppliersAdmin.
	Их метрики надо считать отдельно, чтобы не получить среднюю температуру по больнице.
	Также в рамках каждого клиента желательно раздельно учитывать разные ручки.
	Хотя бы потому что 400/500 ошибки часто специфичны для них.
*/
func NewHttpClientMetrics(clientName, requestName string) IHttpClientMetrics {
	labels := map[string]string{
		metricsLabelClient:  clientName,
		metricsLabelRequest: requestName,
	}

	return &httpClientMetrics{
		NbSuccess:         newCounter(metricsNamespace, metricsSubsystemHttpClient, "nb_success", labels),
		NbError:           newCounter(metricsNamespace, metricsSubsystemHttpClient, "nb_error", labels),
		NbBadStatus:       newCounter(metricsNamespace, metricsSubsystemHttpClient, "nb_bad_status", labels),
		NbBadStatus4xx:    newCounter(metricsNamespace, metricsSubsystemHttpClient, "nb_bad_status_4xx", labels),
		NbBadStatus5xx:    newCounter(metricsNamespace, metricsSubsystemHttpClient, "nb_bad_status_5xx", labels),
		RequestDurationMs: newSummary(metricsNamespace, metricsSubsystemHttpClient, "req_duration_ms", labels),
	}
}

type httpClientMetrics struct {
	NbSuccess         prometheus.Counter
	NbError           prometheus.Counter
	NbBadStatus       prometheus.Counter
	NbBadStatus4xx    prometheus.Counter
	NbBadStatus5xx    prometheus.Counter
	RequestDurationMs prometheus.Summary
}

func (m *httpClientMetrics) RequestDuration(duration time.Duration) {
	m.RequestDurationMs.Observe(float64(duration.Milliseconds()))
}
func (m *httpClientMetrics) IncSuccess() {
	m.NbSuccess.Inc()
}
func (m *httpClientMetrics) IncError() {
	m.NbError.Inc()
}
func (m *httpClientMetrics) IncBadStatus(status int) {
	m.NbBadStatus.Inc()
	if status >= 400 && status <= 499 {
		m.NbBadStatus4xx.Inc()
	} else if status >= 500 && status <= 599 {
		m.NbBadStatus5xx.Inc()
	}
}
