package updater

import "gitlab-metrics/gitlab/dto"

type UserEventsUpdater interface {
	UpdateMetrics(mrs <-chan dto.UserEvent)
	ResetMetrics()
}
