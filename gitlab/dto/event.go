package dto

import "time"

type Event struct {
	Title          string    `json:"title" db:"title"`
	ProjectId      uint32    `json:"project_id" db:"project_id"`
	ActionName     string    `json:"action_name" db:"action_name"`
	TargetId       uint32    `json:"target_id" db:"target_id"`
	TargetType     string    `json:"target_type" db:"target_type"`
	AuthorId       uint32    `json:"author_id" db:"author_id"`
	TargetTitle    string    `json:"target_title" db:"target_title"`
	Author         User      `json:"author" db:"author"`
	AuthorUsername string    `json:"author_username" db:"author_username"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

const (
	EventActionNameAccepted  = "accepted"
	EventActionNameClosed    = "closed"
	EventActionNameCommented = "commented on"
	EventActionNameCreated   = "created"
	EventActionNameDeleted   = "deleted"
	EventActionNameMerged    = "merged"
	EventActionNameOpened    = "opened"
	EventActionNamePushedNew = "pushed new"
	EventActionNamePushedTo  = "pushed to"
	EventActionNameReopened  = "reopened"
	EventActionNameUpdated   = "updated"

	EventTargetTypeDiffNote       = "DiffNote"
	EventTargetTypeDiscussionNote = "DiscussionNote"
	EventTargetTypeIssue          = "Issue"
	EventTargetTypeMergeRequest   = "MergeRequest"
	EventTargetTypeNote           = "Note"
)
