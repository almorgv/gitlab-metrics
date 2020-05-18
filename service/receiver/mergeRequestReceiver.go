package receiver

import (
	"fmt"
	"sync/atomic"
	"time"

	"gitlab-metrics/gitlab/api"
	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

type MergeRequestReceiver interface {
	FetchNew(mrs chan<- dto.ProjectMergeRequest) (uint32, error)
	FetchNewForProject(projectId uint32, mrs chan<- dto.ProjectMergeRequest) (uint32, error)
	FetchNewForAllProjects(mrs chan<- dto.ProjectMergeRequest) (uint32, error)
}

type mergeRequestReceiver struct {
	log.Loggable
	gitlabClient  api.Client
	lastFetchTime time.Time
}

type fetchFunc func(api.ProjectMergeRequestOpts) ([]dto.MergeRequest, error)

func NewMergeRequestReceiver(gitlabClient api.Client) *mergeRequestReceiver {
	return &mergeRequestReceiver{gitlabClient: gitlabClient}
}

func (r *mergeRequestReceiver) FetchNew(mrs chan<- dto.ProjectMergeRequest) (uint32, error) {
	currentFetchTime := time.Now()
	opts := api.ProjectMergeRequestOpts{
		State:         "all",
		UpdatedBefore: currentFetchTime,
		UpdatedAfter:  r.lastFetchTime,
	}

	mrsChan := make(chan dto.MergeRequest)
	go r.mergeProjectAndMrs(dto.Project{}, mrsChan, mrs)

	count, err := r.fetchNewMergeRequestsWith(mrsChan, r.makeFetcher(), opts)
	close(mrsChan)

	if err == nil {
		r.lastFetchTime = currentFetchTime
	}
	return count, err
}

func (r *mergeRequestReceiver) FetchNewForProject(projectId uint32, mrs chan<- dto.ProjectMergeRequest) (uint32, error) {
	currentFetchTime := time.Now()
	opts := api.ProjectMergeRequestOpts{
		State:         "all",
		UpdatedBefore: currentFetchTime,
		UpdatedAfter:  r.lastFetchTime,
	}

	mrsChan := make(chan dto.MergeRequest)
	go r.mergeProjectAndMrs(dto.Project{}, mrsChan, mrs)

	count, err := r.fetchNewMergeRequestsWith(mrsChan, r.makeProjectFetcher(projectId), opts)
	close(mrsChan)

	if err == nil {
		r.lastFetchTime = currentFetchTime
	}
	return count, err
}

func (r *mergeRequestReceiver) FetchNewForAllProjects(mrs chan<- dto.ProjectMergeRequest) (uint32, error) {
	return r.getMergeRequestsForAllProjects(mrs)
}

func (r *mergeRequestReceiver) makeFetcher() fetchFunc {
	return r.gitlabClient.GetMergeRequestsWithOpts
}

func (r *mergeRequestReceiver) makeProjectFetcher(projectId uint32) fetchFunc {
	return func(opts api.ProjectMergeRequestOpts) ([]dto.MergeRequest, error) {
		return r.gitlabClient.GetProjectMergeRequestsWithOpts(projectId, opts)
	}
}

func (r *mergeRequestReceiver) makeAllProjectsFetcher(projectId uint32) fetchFunc {
	return func(opts api.ProjectMergeRequestOpts) ([]dto.MergeRequest, error) {
		return r.gitlabClient.GetProjectMergeRequestsWithOpts(projectId, opts)
	}
}

func (r *mergeRequestReceiver) getMergeRequestsForAllProjects(mrs chan<- dto.ProjectMergeRequest) (uint32, error) {
	totalFetchedProjectsCount := uint32(0)
	totalFetchedMrsCount := uint32(0)
	projectsChan := make(chan dto.Project)

	r.Log().Infof("staring to fetch merge requests for all projects since %s", r.lastFetchTime)
	go func() {
		defer close(projectsChan)
		fetchedProjectsCount, err := r.gitlabClient.GetAllProjectsChanneled(projectsChan)
		atomic.AddUint32(&totalFetchedProjectsCount, fetchedProjectsCount)
		r.Log().Debugf("fetched total %d projects", totalFetchedProjectsCount)
		if err != nil {
			r.Log().ErrorWithErr(err, "An error occurred while fetching projects")
		}
	}()

	currentFetchTime := time.Now()
	opts := api.ProjectMergeRequestOpts{
		State:         "all",
		UpdatedBefore: currentFetchTime,
		UpdatedAfter:  r.lastFetchTime,
	}

	for {
		project, ok := <-projectsChan
		if !ok {
			r.Log().Infof("fetched total %d merge requests over %d projects since %s", totalFetchedMrsCount, totalFetchedProjectsCount, r.lastFetchTime)
			r.lastFetchTime = currentFetchTime
			return totalFetchedMrsCount, nil
		}

		r.Log().Debugf("starting to fetch merge requests since %s for project %d %s", r.lastFetchTime, project.Id, project.PathWithNamespace)

		mrsChan := make(chan dto.MergeRequest)
		go r.mergeProjectAndMrs(project, mrsChan, mrs)

		fetchedCount, err := r.fetchNewMergeRequestsWith(mrsChan, r.makeProjectFetcher(project.Id), opts)
		r.Log().Debugf("fetched %d merge requests for project %d %s", fetchedCount, project.Id, project.PathWithNamespace)

		close(mrsChan)
		totalFetchedMrsCount += fetchedCount

		if err != nil {
			return totalFetchedMrsCount, err
		}
	}
}

func (r *mergeRequestReceiver) mergeProjectAndMrs(project dto.Project, mrs <-chan dto.MergeRequest, projMrs chan<- dto.ProjectMergeRequest) {
	for {
		mr, ok := <-mrs
		if !ok {
			return
		}
		projMrs <- dto.ProjectMergeRequest{
			Mr:      mr,
			Project: project,
		}
	}
}

func (r *mergeRequestReceiver) fetchNewMergeRequestsWith(mrs chan<- dto.MergeRequest, fetchFunc fetchFunc, opts api.ProjectMergeRequestOpts) (uint32, error) {
	fetchedCount := uint32(0)
	pageNumber := uint32(1)

	opts.PerPage = 100

	for {
		opts.Page = pageNumber
		r.Log().Debug("start fetching merge requests")
		mergeRequests, err := fetchFunc(opts)
		r.Log().Debug("end fetching merge requests")
		if err != nil {
			return fetchedCount, fmt.Errorf("fetch new merge requests: %v", err)
		}

		r.Log().Debugf("fetched %d new MRs on page %d", len(mergeRequests), pageNumber)
		if len(mergeRequests) == 0 {
			return fetchedCount, nil
		}

		for _, mr := range mergeRequests {
			mrs <- mr
			fetchedCount++
		}

		pageNumber++
	}
}
