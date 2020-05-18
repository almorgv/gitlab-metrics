package prometheus

import (
	"strconv"
	"time"

	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/metric/definition"
	"gitlab-metrics/metric/updater"
)

type projectEventsUpdater struct {
}

func NewProjectEventsUpdater() *projectEventsUpdater {
	return &projectEventsUpdater{}
}

func (c *projectEventsUpdater) UpdateMetrics(projectEvents <-chan dto.ProjectEvent) {
	for {
		projectEvent, ok := <-projectEvents
		if !ok {
			break
		}
		if isInterestedEvent(projectEvent.Event) {
			c.UpdateMetricsFromEvent(projectEvent.Event)
		}
	}
}

func (c *projectEventsUpdater) ResetMetrics() {
	definition.EventsTotalCount.Reset()
	definition.EventsActions.Reset()
}

func (c *projectEventsUpdater) UpdateMetricsFromEvent(event dto.Event) {
	definition.EventsTotalCount.WithLabelValues(c.getTotalCountLabels(event)).Inc()
	definition.EventsActions.WithLabelValues(c.getActionLabels(event)).Inc()
}

func (c *projectEventsUpdater) getTotalCountLabels(event dto.Event) string {
	return event.AuthorUsername
}

func (c *projectEventsUpdater) getActionLabels(event dto.Event) (string, string, string, string) {
	return event.AuthorUsername, event.ActionName, event.TargetType, uint32ToStr(event.ProjectId)
}

func uint32ToStr(val uint32) string {
	return strconv.FormatUint(uint64(val), 10)
}

func isInterestedEvent(event dto.Event) bool {
	// skip events created not today
	if event.CreatedAt.UTC().Truncate(24*time.Hour) != time.Now().UTC().Truncate(24*time.Hour) {
		return false
	}

	// skip events with not interested actions
	for _, val := range updater.InterestedActions {
		if val == event.ActionName {
			return true
		}
	}

	return false
}
