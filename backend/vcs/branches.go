package vcs

import (
	"github.com/google/uuid"
)

type Branch struct {
	Name        string
	Commits     []uuid.UUID
	CommitOrder map[uuid.UUID]int
}

type BranchHolder struct {
	Branches map[string]Branch
}

var branchHolderInstance *BranchHolder

func GetBranchHolder() *BranchHolder {
	if branchHolderInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if branchHolderInstance == nil {
			branchHolderInstance = &BranchHolder{}
		}
	}
	return branchHolderInstance
}

func CreateNewBranch(name string, currentBranch string) {
	var branches = GetBranchHolder()
	val, currOk := branches.branches[currentBranch]
	if currOk {
		branches.branches[name] = *val.Clone(name)
	}
}

func (b *Branch) GetName() string              { return b.name }
func (b *Branch) GetCommit(idx uint) uuid.UUID { return b.commits[idx] }
func (b *Branch) GetCommitsRange(from uuid.UUID, to uuid.UUID) []uuid.UUID {
	var ret = make([]uuid.UUID, 0)

	if from == to {
		return append(ret, from)
	}

	var fromNum = b.commitOrder[from]
	var toNum = b.commitOrder[to]
	var first = min(fromNum, toNum)
	var last = max(fromNum, toNum)
	return b.commits[first : last+1]
}
func (b *Branch) Clone(newName string) *Branch {
	var br = &Branch{}
	br.name = newName
	copy(br.commits, b.commits)
	for k, v := range b.commitOrder {
		br.commitOrder[k] = v
	}
	return br
}
func (b *Branch) GetDiffsInBranch() []Diff {
	var diffs = make([]Diff, 0)
	for _, v := range b.commits {
		diffs = append(diffs, GetDiff(v))
	}
	return diffs
}
