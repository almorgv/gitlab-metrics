package service

import (
	"github.com/jasonlvhit/gocron"

	"gitlab-metrics/log"
)

type Scheduler interface {
	Submit() error
	Run()
}

type scheduler struct {
	instance *gocron.Scheduler
	log.Loggable
}

func NewScheduler() *scheduler {
	return &scheduler{
		//instance: gocron.NewScheduler(),
	}
}

func (s *scheduler) Submit(interval uint64, job interface{}, params ...interface{}) error {
	return gocron.Every(interval).Minute().From(gocron.NextTick()).Do(job, params...)
}

func (s *scheduler) Run() {
	<-gocron.Start()
}
