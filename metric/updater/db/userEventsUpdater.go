package db

import (
	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
	"gitlab-metrics/metric/updater"
	"gitlab-metrics/repository"
)

type userEventsUpdater struct {
	eventRepository repository.EventRepository
	log.Loggable
}

func NewUserEventsUpdater(eventRepository repository.EventRepository) *userEventsUpdater {
	return &userEventsUpdater{
		eventRepository: eventRepository,
	}
}

func (u *userEventsUpdater) ResetMetrics() {

}

func (u *userEventsUpdater) UpdateMetrics(userEvents <-chan dto.UserEvent) {
	for {
		userEvent, ok := <-userEvents
		if !ok {
			break
		}
		if !u.isInterestedEvent(userEvent.Event) {
			continue
		}
		if err := u.eventRepository.Create(userEvent.Event); err != nil {
			u.Log().Errorf("Failed to create event: %v", err)
		}

	}
}

func (u *userEventsUpdater) isInterestedEvent(event dto.Event) bool {
	for _, val := range updater.InterestedActions {
		if val == event.ActionName {
			return true
		}
	}
	return false
}
