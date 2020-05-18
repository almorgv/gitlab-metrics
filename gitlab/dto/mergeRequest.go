package dto

import "time"

type MergeRequest struct {
	BlockingDiscussionsResolved bool      `json:"blocking_discussions_resolved"`
	ChangesCount                string    `json:"changes_count"`
	CreatedAt                   time.Time `json:"created_at"`
	Description                 string    `json:"description"`
	DiscussionLocked            bool      `json:"discussion_locked"`
	DivergedCommitsCount        uint16    `json:"diverged_commits_count"`
	Downvotes                   uint16    `json:"downvotes"`
	HasConflicts                bool      `json:"has_conflicts"`
	Id                          uint32    `json:"id"`
	Iid                         uint32    `json:"iid"`
	Labels                      []string  `json:"labels"`
	LatestBuildFinishedAt       time.Time `json:"latest_build_finished_at"`
	LatestBuildStartedAt        time.Time `json:"latest_build_started_at"`
	MergedAt                    time.Time `json:"merged_at"`
	MergeError                  string    `json:"merge_error"`
	MergeStatus                 string    `json:"merge_status"`
	ProjectId                   uint32    `json:"project_id"`
	SourceBranch                string    `json:"source_branch"`
	State                       string    `json:"state"`
	TargetBranch                string    `json:"target_branch"`
	Title                       string    `json:"title"`
	UpdatedAt                   time.Time `json:"updated_at"`
	Upvotes                     uint16    `json:"upvotes"`
	UserNotesCount              uint16    `json:"user_notes_count"`
	WebUrl                      string    `json:"web_url"`
	WorkInProgress              bool      `json:"work_in_progress"`
	Author                      User      `json:"author"`
	MergedBy                    User      `json:"merged_by"`
	Project                     Project
}

const (
	MergeRequestStateOpened = "opened"
	MergeRequestStateMerged = "merged"
	MergeRequestStateClosed = "closed"
)
