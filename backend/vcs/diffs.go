package vcs

import (
	"github.com/google/uuid"
	"image"
	"log"
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
	Diffs map[uuid.UUID]*Diff
}

var commitHolderInstance *CommitHolder

func GetCommitHolder() *CommitHolder {
	if commitHolderInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if commitHolderInstance == nil {
			commitHolderInstance = &CommitHolder{make(map[uuid.UUID]*Diff)}
		}
	}
	return commitHolderInstance
}

func GetDiff(commit uuid.UUID) *Diff {
	var diffList = GetCommitHolder()
	return diffList.Diffs[commit]
}

func CreateInitialDiff(image image.RGBA) []PixelDiff {

	bounds := image.Bounds()
	xMax := bounds.Max.X
	xMin := bounds.Min.X
	yMax := bounds.Max.Y
	yMin := bounds.Min.Y

	log.Printf("New image bounds: %d, %d, %d, %d", xMin, xMax, yMin, yMax)

	diffs := make([]PixelDiff, (xMax-xMin)*(yMax-xMin))

	lazy := 0

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			r, g, b, a := image.At(x, y).RGBA()
			diffs[lazy] = PixelDiff{int16(x), int16(y), int16(r), int16(g), int16(b), int16(a)}
			lazy++
		}
	}

	return diffs
}

func GetImageDiff(oldImage image.RGBA, newImage image.RGBA) []PixelDiff {

	// sorry but I'm not checking the images have the same dimensions
	bounds := oldImage.Bounds()
	xMax := bounds.Max.X
	xMin := bounds.Min.X
	yMax := bounds.Max.Y
	yMin := bounds.Min.Y

	diffs := make([]PixelDiff, (xMax-xMin)*(yMax-xMin))

	lazy := 0

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			r1, g1, b1, a1 := oldImage.At(x, y).RGBA()
			r2, g2, b2, a2 := newImage.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				diffs[lazy] = PixelDiff{int16(x), int16(y), int16(r2), int16(g2), int16(b2), int16(a2)}
				lazy++
			}

		}
	}

	return diffs
}

func CreateCommit(changes []PixelDiff) uuid.UUID {
	var commitId = uuid.New()
	var commits = GetCommitHolder()
	var newCommit = Diff{
		Commit:       commitId,
		PixelChanges: changes,
		Timestamp:    time.Now(),
	}
	commits.Diffs[commitId] = &newCommit
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
