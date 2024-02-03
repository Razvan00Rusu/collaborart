package vcs

import (
	"github.com/google/uuid"
	"sync"
)

type PixelDiff struct {
	X, Y, DR, DG, DB, DA int16
}

type Diff struct {
	Commit       uuid.UUID
	PixelChanges []PixelDiff
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
			commitHolderInstance = &CommitHolder{}
		}
	}
	return commitHolderInstance
}

func GetDiff(commit uuid.UUID) Diff {
	var diffList = GetCommitHolder()
	return diffList.Diffs[commit]
}
