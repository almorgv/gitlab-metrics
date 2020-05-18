package definition

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricsSubsystemMergeRequests = "merge_requests"
)

var (
	mergeRequestMetricLabels = []string{"projectid", "projectpath"}
)

var (
	TotalMergeRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemMergeRequests,
		Name:      "total_count",
		Help:      "The number of total merge requests",
	}, mergeRequestMetricLabels)

	OpenedMergeRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemMergeRequests,
		Name:      "opened",
		Help:      "The number of opened merge requests",
	}, mergeRequestMetricLabels)

	MergedMergeRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemMergeRequests,
		Name:      "merged",
		Help:      "The number of merged merge requests",
	}, mergeRequestMetricLabels)

	ClosedMergeRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemMergeRequests,
		Name:      "closed",
		Help:      "The number of closed merge requests",
	}, mergeRequestMetricLabels)

	MergeRequestUpvotes = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemMergeRequests,
		Name:      "upvotes",
		Help:      "The number of upvotes in merge requests",
	}, mergeRequestMetricLabels)

	MergeRequestNotes = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemMergeRequests,
		Name:      "notes",
		Help:      "The number of notes in merge requests",
	}, mergeRequestMetricLabels)

	MergeRequestChangesCount = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystemMergeRequests,
		Name:      "changes_count",
		Help:      "The number of changes in merge requests",
	}, mergeRequestMetricLabels)
)
