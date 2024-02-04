package vcs

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type PixelDiff struct {
	X, Y, R, G, B, A int16
}

type Diff struct {
	Commit       uuid.UUID
	PixelChanges []PixelDiff
	Timestamp    time.Time
}

var lock = &sync.Mutex{}

type CommitHolder struct {
	Diffs map[uuid.UUID]Diff
}

var commitHolderInstance *CommitHolder

func GetCommitHolder() *CommitHolder {
	if commitHolderInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if commitHolderInstance == nil {
			commitHolderInstance = &CommitHolder{make(map[uuid.UUID]Diff)}
		}
	}
	return commitHolderInstance
}

func GetDiff(commit uuid.UUID) Diff {
	var diffList = GetCommitHolder()
	return diffList.Diffs[commit]
}

func CreateCommit(changes []PixelDiff) uuid.UUID {
	var commitId = uuid.New()
	var commits = GetCommitHolder()
	var newCommit = Diff{
		Commit:       commitId,
		PixelChanges: changes,
		Timestamp:    time.Now(),
	}
	commits.Diffs[commitId] = newCommit
	return commitId
}

func SquashCommitsToPixelChanges(commits []uuid.UUID) []PixelDiff {
	var pixelDiffs = make([]PixelDiff, 0)
	for _, v := range commits {
		var diff = GetDiff(v)
		pixelDiffs = append(pixelDiffs, diff.PixelChanges...)
	}
	return pixelDiffs
}
