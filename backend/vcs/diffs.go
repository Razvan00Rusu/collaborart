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

type commitHolder struct {
	diffs map[uuid.UUID]PixelDiff
}
