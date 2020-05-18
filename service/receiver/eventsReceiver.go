package receiver

import (
	"time"

	"gitlab-metrics/gitlab/api"
	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

type EventReceiver interface {
	FetchNewForAllUsers(eventsChan chan<- dto.UserEvent, since time.Time) (uint32, error)
	FetchNewForAllProjects(eventsChan chan<- dto.ProjectEvent, since time.Time) (uint32, error)
}

type eventReceiver struct {
	log.Loggable
	gitlabClient api.Client
}

func NewEventReceiver(gitlabClient api.Client) *eventReceiver {
	return &eventReceiver{
		gitlabClient: gitlabClient,
	}
}

func (r *eventReceiver) FetchNewForAllUsers(eventsChan chan<- dto.UserEvent, since time.Time) (uint32, error) {
	opts := api.EventOpts{
		After: since,
	}
	return r.gitlabClient.GetEventsForAllUsersChanneledWithOpts(eventsChan, opts)
}

func (r *eventReceiver) FetchNewForAllProjects(eventsChan chan<- dto.ProjectEvent, since time.Time) (uint32, error) {
	opts := api.EventOpts{
		After: since,
	}
	return r.gitlabClient.GetEventsForAllProjectsChanneledWithOpts(eventsChan, opts)
}
