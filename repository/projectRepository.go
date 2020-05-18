package repository

import (
	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

type ProjectRepository interface {
	CreateOrUpdateLastActivity(project dto.Project) error
}

type projectRepository struct {
	db *Db
	log.Loggable
}

func NewProjectRepository(db *Db) *projectRepository {
	return &projectRepository{db: db}
}

func (e *projectRepository) CreateOrUpdateLastActivity(project dto.Project) error {
	_, err := e.db.NamedExec(`insert into 
    	projects(id, name, name_with_namespace, created_at, last_activity_at)
		values (:id, :name, :name_with_namespace, :created_at, :last_activity_at)
		on conflict (id) DO UPDATE SET last_activity_at=:last_activity_at`,
		project)
	return err
}
