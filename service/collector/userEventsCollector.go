package collector

import (
	"time"

	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
	"gitlab-metrics/metric/updater"
	"gitlab-metrics/service/receiver"
)

type UserEventsCollector interface {
	FetchAndUpdateMetrics()
	RegisterUpdater(eventUpdater updater.UserEventsUpdater)
}

type userEventsCollector struct {
	log.Loggable
	eventReceiver  receiver.EventReceiver
	eventUpdaters  map[updater.UserEventsUpdater]chan dto.UserEvent
	lastUpdateTime time.Time
}

func NewUserEventsCollector(eventReceiver receiver.EventReceiver) *userEventsCollector {
	return &userEventsCollector{
		eventReceiver: eventReceiver,
		eventUpdaters: make(map[updater.UserEventsUpdater]chan dto.UserEvent),
	}
}

func (c *userEventsCollector) RegisterUpdater(eventUpdater updater.UserEventsUpdater) {
	// this channel will never be closed
	ch := make(chan dto.UserEvent)
	c.eventUpdaters[eventUpdater] = ch
	c.Log().Debug("Starting event metrics updater")
	go eventUpdater.UpdateMetrics(ch)
}

func (c *userEventsCollector) multiplexOutputs(input <-chan dto.UserEvent) {
	for {
		v, ok := <-input
		if !ok {
			break
		}
		for _, outputChannel := range c.eventUpdaters {
			go func(och chan<- dto.UserEvent) { och <- v }(outputChannel)
		}
	}
}

func (c *userEventsCollector) FetchAndUpdateMetrics() {
	c.Log().Info("Start updating user metrics")

	for updter := range c.eventUpdaters {
		c.Log().Debug("Resetting event metrics")
		updter.ResetMetrics()
	}

	events := make(chan dto.UserEvent)
	defer close(events)

	go c.multiplexOutputs(events)

	count, err := c.eventReceiver.FetchNewForAllUsers(events, c.getTimeToFetchEventsAfter())
	if err != nil {
		c.Log().ErrorfWithErr(err, "An error occurred while fetching new user events: %v", err)
	}
	c.Log().Infof("Metrics updated with %d new user events", count)
	c.lastUpdateTime = time.Now()
}

// Returns zero date if events were never fetched before to be able to get all events up to now.
// Otherwise returns the yesterday date with the time set to 23:59:59.
// Because if `after` option is specified in the API request
// then only events with created date greater then specified will be fetched.
// To get events for today this option must be set to the yesterday date.
// Time is added only for clarification purpose and means that we dont want to get events for yesterday.
func (c *userEventsCollector) getTimeToFetchEventsAfter() time.Time {
	// Don't set `after` option if events were never fetched before to be able to get all events up to now
	if c.lastUpdateTime.IsZero() {
		return time.Time{}
	}
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, time.UTC)
}
