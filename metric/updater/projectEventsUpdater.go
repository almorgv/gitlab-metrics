package updater

import "gitlab-metrics/gitlab/dto"

type ProjectEventsUpdater interface {
	UpdateMetrics(mrs <-chan dto.ProjectEvent)
	ResetMetrics()
}
