package bsky

type CommitEvent interface {
	Type() CommitEventType
	GetRepo() string
	GetOps() []RepoOperation
	GetBlocks() []byte
}

type CommitEventType string

const (
	CommitEventTypeRepoCommit CommitEventType = "repo_commit"
)
