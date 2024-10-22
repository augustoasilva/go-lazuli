package dto

type CommitEvent interface {
	Type() CommitEventType
}

type CommitEventType string

const (
	CommitEventTypeRepoCommit CommitEventType = "repo_commit"
)
