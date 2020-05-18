package api

import (
	"fmt"
	"path"
	"sync/atomic"

	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

const (
	eventsPath = "/events"
)

type EventsApi interface {
	GetUserEventsWithOpts(userId uint32, opts EventOpts) ([]dto.Event, error)
	GetEventsForAllUsersChanneledWithOpts(userEvents chan<- dto.UserEvent, opts EventOpts) (uint32, error)

	GetProjectEventsWithOpts(projectId uint32, opts EventOpts) ([]dto.Event, error)
	GetEventsForAllProjectsChanneledWithOpts(projectEvents chan<- dto.ProjectEvent, opts EventOpts) (uint32, error)
}

type eventsApi struct {
	api
	projectsApi
	usersApi
	log.Loggable
}

func (e *eventsApi) GetUserEventsWithOpts(userId uint32, opts EventOpts) ([]dto.Event, error) {
	urlPath := e.makeUserEventsUrlPath(userId)
	return e.getEventsWithOpts(urlPath, opts)
}

func (e *eventsApi) GetProjectEventsWithOpts(projectId uint32, opts EventOpts) ([]dto.Event, error) {
	urlPath := e.makeProjectEventsUrlPath(projectId)
	return e.getEventsWithOpts(urlPath, opts)
}

func (e *eventsApi) GetEventsForAllUsersChanneledWithOpts(userEvents chan<- dto.UserEvent, opts EventOpts) (uint32, error) {
	totalFetchedUsersCount := uint32(0)
	totalFetchedEventsCount := uint32(0)
	users := make(chan dto.User)

	e.Log().Infof("staring to fetch events for all users with %s", opts)
	go func() {
		defer close(users)
		fetchedUsersCount, err := e.GetAllUsersChanneled(users)
		atomic.AddUint32(&totalFetchedUsersCount, fetchedUsersCount)
		e.Log().Debugf("fetched total %d users", totalFetchedUsersCount)
		if err != nil {
			e.Log().ErrorWithErr(err, "An error occurred while fetching users")
		}
	}()

	events := make(chan dto.Event)
	defer close(events)

	for {
		user, ok := <-users
		if !ok {
			e.Log().Infof("fetched total %d events over %d users with %s", totalFetchedEventsCount, totalFetchedUsersCount, opts)
			return totalFetchedEventsCount, nil
		}

		e.Log().Debugf("starting to fetch events for user %d %s with", user.Id, user.Username, opts)

		go e.mergeUserAndEvents(user, events, userEvents)

		fetchedCount, err := e.getEventsPaged(e.makeUserEventsUrlPath(user.Id), opts, events)
		e.Log().Debugf("fetched %d events for user %d %s", fetchedCount, user.Id, user.Username)

		totalFetchedEventsCount += fetchedCount

		if err != nil {
			return totalFetchedEventsCount, err
		}
	}
}

func (e *eventsApi) GetEventsForAllProjectsChanneledWithOpts(projectEvents chan<- dto.ProjectEvent, opts EventOpts) (uint32, error) {
	totalFetchedProjectsCount := uint32(0)
	totalFetchedEventsCount := uint32(0)
	projects := make(chan dto.Project)

	e.Log().Infof("staring to fetch events for all projects with %s", opts)
	go func() {
		defer close(projects)
		fetchedUsersCount, err := e.GetAllProjectsChanneled(projects)
		atomic.AddUint32(&totalFetchedProjectsCount, fetchedUsersCount)
		e.Log().Debugf("fetched total %d projects", totalFetchedProjectsCount)
		if err != nil {
			e.Log().ErrorWithErr(err, "An error occurred while fetching projects")
		}
	}()

	for {
		project, ok := <-projects
		if !ok {
			e.Log().Infof("fetched total %d events over %d projects with %s", totalFetchedEventsCount, totalFetchedProjectsCount, opts)
			return totalFetchedEventsCount, nil
		}

		e.Log().Debugf("starting to fetch events for project %d %s with", project.Id, project.PathWithNamespace, opts)

		fetchedCount, err := e.GetEventsForProjectWithOpts(project, projectEvents, opts)
		e.Log().Debugf("fetched %d events for project %d %s", fetchedCount, project.Id, project.PathWithNamespace)
		totalFetchedEventsCount += fetchedCount

		if err != nil {
			e.Log().Warnf("An error occurred fetching events for project %d %s: %v", project.Id, project.PathWithNamespace, err)
		}
	}
}

func (e *eventsApi) GetEventsForProjectWithOpts(project dto.Project, projectEvents chan<- dto.ProjectEvent, opts EventOpts) (uint32, error) {
	events := make(chan dto.Event)
	defer close(events)

	go e.mergeProjectsAndEvents(project, events, projectEvents)

	return e.getEventsPaged(e.makeProjectEventsUrlPath(project.Id), opts, events)
}

func (e *eventsApi) makeUserEventsUrlPath(userId uint32) string {
	return fmt.Sprintf(path.Join(apiPath, userPath, eventsPath), userId)
}

func (e *eventsApi) makeProjectEventsUrlPath(projectId uint32) string {
	return fmt.Sprintf(path.Join(apiPath, projectPath, eventsPath), projectId)
}

func (e *eventsApi) getEventsWithOpts(urlPath string, opts EventOpts) ([]dto.Event, error) {
	urlValues := opts.ToValues()
	reqUrl := fmt.Sprintf("%s%s?%s", e.GetBaseUrl(), urlPath, urlValues.Encode())

	var events []dto.Event

	err := e.FetchData(reqUrl, &events)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %v", err)
	}

	return events, nil
}

func (e *eventsApi) getEventsPaged(urlPath string, opts EventOpts, eventsChan chan<- dto.Event) (uint32, error) {
	fetchedCount := uint32(0)
	pageNumber := uint32(1)

	opts.PerPage = 100

	for {
		opts.Page = pageNumber
		e.Log().Tracef("start fetching events by %s", urlPath)
		events, err := e.getEventsWithOpts(urlPath, opts)
		e.Log().Tracef("end fetching events by %s", urlPath)
		if err != nil {
			return fetchedCount, fmt.Errorf("fetch new events: %v", err)
		}

		e.Log().Tracef("fetched %d new events on page %d by %s", len(events), pageNumber, urlPath)
		if len(events) == 0 {
			return fetchedCount, nil
		}

		for _, mr := range events {
			eventsChan <- mr
			fetchedCount++
		}

		pageNumber++
	}
}

func (e *eventsApi) mergeUserAndEvents(user dto.User, events <-chan dto.Event, userEvents chan<- dto.UserEvent) {
	for {
		event, ok := <-events
		if !ok {
			return
		}
		userEvents <- dto.UserEvent{
			User:  user,
			Event: event,
		}
	}
}

func (e *eventsApi) mergeProjectsAndEvents(project dto.Project, events <-chan dto.Event, projectEvents chan<- dto.ProjectEvent) {
	for {
		event, ok := <-events
		if !ok {
			return
		}
		projectEvents <- dto.ProjectEvent{
			Project: project,
			Event:   event,
		}
	}
}
