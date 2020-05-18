package dto

import "time"

type Project struct {
	Id                uint32     `json:"id" db:"id"`
	WebUrl            string     `json:"web_url" db:"web_url"`
	Name              string     `json:"name" db:"name"`
	NameWithNamespace string     `json:"name_with_namespace" db:"name_with_namespace"`
	Path              string     `json:"path" db:"path"`
	PathWithNamespace string     `json:"path_with_namespace" db:"path_with_namespace"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	LastActivityAt    time.Time  `json:"last_activity_at" db:"last_activity_at"`
	Statistics        Statistics `json:"statistics" db:"statistics"`
}

type Statistics struct {
	CommitCount      uint64 `json:"commit_count"`
	StorageSize      uint64 `json:"storage_size"`
	RepositorySize   uint64 `json:"repository_size"`
	WikiSize         uint64 `json:"wiki_size"`
	LfsObjectsSize   uint64 `json:"lfs_objects_size"`
	JobArtifactsSize uint64 `json:"job_artifacts_size"`
	PackagesSize     uint64 `json:"packages_size"`
}
