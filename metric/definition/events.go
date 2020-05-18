package definition

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricsSubsystemEvents = "user_events"
)

var (
	eventLabelsUsername         = []string{"username"}
	eventLabelsUserActionTarget = []string{"username", "action", "target", "projectid"}
)

var (
	EventsTotalCount = makeGauge(
		"total_count",
		"Total count of user events",
		eventLabelsUsername,
	)

	EventsActions = makeGauge(
		"actions",
		"User actions",
		eventLabelsUserActionTarget,
	)
)

func makeGauge(name string, help string, labels []string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemEvents,
		Name:      name,
		Help:      help,
	}, labels)
}
