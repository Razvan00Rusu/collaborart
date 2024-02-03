package vcs

import (
	"github.com/google/uuid"
	"sync"
)

type PixelDiff struct {
	x, y, dR, dG, dB, dA int16
}

type Diff struct {
	commit       uuid.UUID
	pixelChanges []PixelDiff
}

var lock = &sync.Mutex{}

type CommitHolder struct {
	diffs map[uuid.UUID]Diff
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
	return diffList.diffs[commit]
}
