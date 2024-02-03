package vcs

import (
	"github.com/google/uuid"
)

type Branch struct {
	name        string
	commits     []uuid.UUID
	commitOrder map[uuid.UUID]int
}

func (b *Branch) GetName() string              { return b.name }
func (b *Branch) GetCommit(idx uint) uuid.UUID { return b.commits[idx] }
func (b *Branch) GetCommitsRange(from uuid.UUID, to uuid.UUID) []uuid.UUID {
	var ret []uuid.UUID = make([]uuid.UUID, 0)

	if from == to {
		return append(ret, from)
	}

	var fromNum = b.commitOrder[from]
	var toNum = b.commitOrder[to]
	var first = min(fromNum, toNum)
	var last = max(fromNum, toNum)
	return b.commits[first : last+1]
}
