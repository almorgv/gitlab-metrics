package updater

import "gitlab-metrics/gitlab/dto"

type MergeRequestUpdater interface {
	UpdateMetrics(mrs <-chan dto.ProjectMergeRequest)
}
