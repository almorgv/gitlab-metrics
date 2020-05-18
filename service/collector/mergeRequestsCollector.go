package collector

import (
	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
	"gitlab-metrics/metric/updater"
	"gitlab-metrics/service/receiver"
)

type MergeRequestsCollector interface {
	FetchAndUpdateMetrics()
}

type mergeRequestsCollector struct {
	log.Loggable
	mergeRequestReceiver receiver.MergeRequestReceiver
	mergeRequestUpdater  updater.MergeRequestUpdater
}

func NewMergeRequestsCollector(mergeRequestReceiver receiver.MergeRequestReceiver, mergeRequestUpdater updater.MergeRequestUpdater) *mergeRequestsCollector {
	return &mergeRequestsCollector{
		mergeRequestReceiver: mergeRequestReceiver,
		mergeRequestUpdater:  mergeRequestUpdater,
	}
}

func (s *mergeRequestsCollector) FetchAndUpdateMetrics() {
	s.Log().Info("Start updating merge request metrics")
	mergeRequests := make(chan dto.ProjectMergeRequest)
	defer close(mergeRequests)

	s.Log().Debug("Starting merge request metrics updater")
	go s.mergeRequestUpdater.UpdateMetrics(mergeRequests)

	count, err := s.mergeRequestReceiver.FetchNewForAllProjects(mergeRequests)
	if err != nil {
		s.Log().ErrorfWithErr(err, "An error occurred while fetching new merge requests: %v", err)
	}
	s.Log().Infof("Metrics updated with %d new merge requests", count)
}
