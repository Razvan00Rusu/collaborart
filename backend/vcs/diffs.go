package vcs

import (
	"github.com/google/uuid"
	"sync"
)

type pixelDiff struct {
	x, y, dR, dG, dB, dA int16
}

type diff struct {
	commit       uuid.UUID
	pixelChanges []pixelDiff
}

var lock = &sync.Mutex{}

type commitHolder struct {
	diffs map[uuid.UUID]pixelDiff
}
