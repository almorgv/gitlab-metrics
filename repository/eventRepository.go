package repository

import (
	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

type EventRepository interface {
	Create(event dto.Event) error
}

type eventRepository struct {
	db *Db
	log.Loggable
}

func NewEventRepository(db *Db) *eventRepository {
	return &eventRepository{db: db}
}

func (e *eventRepository) Create(event dto.Event) error {
	_, err := e.db.NamedExec(`insert into 
    	events(project_id, action_name, target_id, target_type, author_id, author_username, created_at)
		values (:project_id, :action_name, :target_id, :target_type, :author_id, :author_username, :created_at)
		on conflict DO NOTHING`,
		event)
	return err
}
