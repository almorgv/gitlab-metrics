package prometheus

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/metric/definition"
)

type mergeRequestUpdater struct {
}

func NewMergeRequestUpdater() *mergeRequestUpdater {
	return &mergeRequestUpdater{}
}

func (c *mergeRequestUpdater) UpdateMetrics(mrs <-chan dto.ProjectMergeRequest) {
	for {
		mr, ok := <-mrs
		if !ok {
			break
		}
		c.UpdateMetricsFromMergeRequest(mr)
	}
}

func (c *mergeRequestUpdater) UpdateMetricsFromMergeRequest(projmr dto.ProjectMergeRequest) {
	c.getStateMetric(projmr).Inc()
	c.getPromGauge(definition.TotalMergeRequests, projmr).Inc()
	c.getPromObserver(definition.MergeRequestUpvotes, projmr).Observe(float64(projmr.Mr.Upvotes))
	c.getPromObserver(definition.MergeRequestNotes, projmr).Observe(float64(projmr.Mr.UserNotesCount))
	// FIXME: changes count is always empty
	changesCount, _ := strconv.ParseFloat(strings.TrimRight(projmr.Mr.ChangesCount, "+"), 64)
	c.getPromObserver(definition.MergeRequestChangesCount, projmr).Observe(changesCount)
}

func (c *mergeRequestUpdater) getMetricLabels(projmr dto.ProjectMergeRequest) (string, string) {
	return strconv.FormatUint(uint64(projmr.Mr.ProjectId), 10), projmr.Project.PathWithNamespace
}

func (c *mergeRequestUpdater) getPromGauge(vec *prometheus.GaugeVec, projmr dto.ProjectMergeRequest) prometheus.Gauge {
	return vec.WithLabelValues(c.getMetricLabels(projmr))
}

func (c *mergeRequestUpdater) getPromObserver(vec prometheus.ObserverVec, projmr dto.ProjectMergeRequest) prometheus.Observer {
	return vec.WithLabelValues(c.getMetricLabels(projmr))
}

func (c *mergeRequestUpdater) getStateMetric(projmr dto.ProjectMergeRequest) prometheus.Gauge {
	switch projmr.Mr.State {
	case dto.MergeRequestStateOpened:
		return c.getPromGauge(definition.OpenedMergeRequests, projmr)
	case dto.MergeRequestStateMerged:
		return c.getPromGauge(definition.MergedMergeRequests, projmr)
	case dto.MergeRequestStateClosed:
		return c.getPromGauge(definition.ClosedMergeRequests, projmr)
	default:
		return nil
	}
}
