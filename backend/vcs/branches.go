package vcs

import (
	"github.com/google/uuid"
)

type Branch struct {
	Name          string
	Commits       []uuid.UUID
	CommitOrder   map[uuid.UUID]int
	Width, Height int16
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
	val, currOk := branches.Branches[currentBranch]
	if currOk {
		branches.Branches[name] = *val.Clone(name)
	}
}

func (b *Branch) GetName() string              { return b.Name }
func (b *Branch) GetCommit(idx uint) uuid.UUID { return b.Commits[idx] }
func (b *Branch) GetCommitsRange(from uuid.UUID, to uuid.UUID) []uuid.UUID {
	var ret = make([]uuid.UUID, 0)

	if from == to {
		return append(ret, from)
	}

	var fromNum = b.CommitOrder[from]
	var toNum = b.CommitOrder[to]
	var first = min(fromNum, toNum)
	var last = max(fromNum, toNum)
	return b.Commits[first : last+1]
}
func (b *Branch) AddCommit(changes []PixelDiff) {
	var newCommitId = CreateCommit(changes)
	b.Commits = append(b.Commits, newCommitId)
}
func (b *Branch) Clone(newName string) *Branch {
	var br = &Branch{}
	br.Name = newName
	copy(br.Commits, b.Commits)
	for k, v := range b.CommitOrder {
		br.CommitOrder[k] = v
	}
	return br
}
func (b *Branch) GetDiffsInBranch() []Diff {
	var diffs = make([]Diff, 0)
	for _, v := range b.Commits {
		diffs = append(diffs, GetDiff(v))
	}
	return diffs
}

func BranchExists(name string) bool {
	var branches = GetBranchHolder()
	_, ok := branches.Branches[name]
	return ok
}

func CreateOrphanBranch(name string) {
	if BranchExists(name) {
		panic("Branch already exists!")
	}
	var branches = GetBranchHolder()
	var newCommit = CreateCommit(nil)
	var mainBranch = Branch{
		Name:        name,
		Commits:     append(make([]uuid.UUID, 0), newCommit),
		CommitOrder: map[uuid.UUID]int{newCommit: 0},
		Width:       0,
		Height:      0,
	}
	branches.Branches[name] = mainBranch
}
