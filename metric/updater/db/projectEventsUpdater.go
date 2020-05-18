package db

import (
	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
	"gitlab-metrics/metric/updater"
	"gitlab-metrics/repository"
)

type projectEventsUpdater struct {
	eventRepository   repository.EventRepository
	projectRepository repository.ProjectRepository
	log.Loggable
}

func NewProjectEventsUpdater(
	eventRepository repository.EventRepository,
	projectRepository repository.ProjectRepository,
) *projectEventsUpdater {
	return &projectEventsUpdater{
		eventRepository:   eventRepository,
		projectRepository: projectRepository,
	}
}

func (p *projectEventsUpdater) ResetMetrics() {

}

func (p *projectEventsUpdater) UpdateMetrics(projectEvents <-chan dto.ProjectEvent) {
	for {
		projectEvent, ok := <-projectEvents
		if !ok {
			break
		}
		if !p.isInterestedEvent(projectEvent.Event) {
			continue
		}
		if err := p.projectRepository.CreateOrUpdateLastActivity(projectEvent.Project); err != nil {
			p.Log().Errorf("Failed to create project: %v", err)
		}
		if err := p.eventRepository.Create(projectEvent.Event); err != nil {
			p.Log().Errorf("Failed to create event: %v", err)
		}

	}
}

func (p *projectEventsUpdater) isInterestedEvent(event dto.Event) bool {
	for _, val := range updater.InterestedActions {
		if val == event.ActionName {
			return true
		}
	}
	return false
}
