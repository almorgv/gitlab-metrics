package api

import (
	"fmt"
	"path"

	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

const (
	projectsPath = "/projects"
	projectPath  = "/projects/%d"
)

type ProjectsApi interface {
	GetAllProjects() ([]dto.Project, error)
	GetAllProjectsWithOpts(opts ProjectsOpts) ([]dto.Project, error)
	GetAllProjectsChanneled(fetchedProjects chan<- dto.Project) (uint32, error)

	GetProject(projectId uint32) (*dto.Project, error)
	GetProjectWithOpts(projectId uint32, opts ProjectsOpts) (*dto.Project, error)
}

type projectsApi struct {
	api
	log.Loggable
}

func (p *projectsApi) GetAllProjectsWithOpts(opts ProjectsOpts) ([]dto.Project, error) {
	urlPath := path.Join(apiPath, projectsPath)
	urlValues := opts.ToValues()
	reqUrl := fmt.Sprintf("%s%s?%s", p.GetBaseUrl(), urlPath, urlValues.Encode())

	var projects []dto.Project

	err := p.FetchData(reqUrl, &projects)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %v", err)
	}

	return projects, nil
}

func (p *projectsApi) GetAllProjects() ([]dto.Project, error) {
	return p.GetAllProjectsWithOpts(ProjectsOpts{})
}

func (p *projectsApi) GetAllProjectsChanneled(fetchedProjects chan<- dto.Project) (uint32, error) {
	fetchedCount := uint32(0)
	pageNumber := uint32(1)

	for {
		projects, err := p.GetAllProjectsWithOpts(ProjectsOpts{
			RequestOpts: RequestOpts{
				PerPage: 100,
				Page:    pageNumber,
			}})

		if err != nil {
			return fetchedCount, err
		}

		p.Log().Debugf("fetched %d projects on page %d", len(projects), pageNumber)
		if len(projects) == 0 {
			return fetchedCount, nil
		}

		for _, proj := range projects {
			fetchedProjects <- proj
			fetchedCount++
		}

		pageNumber++
	}
}

func (p *projectsApi) GetProjectWithOpts(projectId uint32, opts ProjectsOpts) (*dto.Project, error) {
	urlPath := fmt.Sprintf(path.Join(apiPath, projectPath), projectId)
	urlValues := opts.ToValues()
	reqUrl := fmt.Sprintf("%s%s?%s", p.GetBaseUrl(), urlPath, urlValues.Encode())

	project := &dto.Project{}

	err := p.FetchData(reqUrl, project)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %v", err)
	}

	return project, nil
}

func (p *projectsApi) GetProject(projectId uint32) (*dto.Project, error) {
	return p.GetProjectWithOpts(projectId, ProjectsOpts{})
}
