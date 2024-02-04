package backend

import (
	"collaborart/backend/composedImage"
	"collaborart/backend/vcs"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/draw"
	"log"
)

func PushToBranch(branchId string, imageFile *image.Image) {

	b := (*imageFile).Bounds()
	imgRGB := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(imgRGB, imgRGB.Bounds(), *imageFile, b.Min, draw.Src)

	if vcs.BranchExists(branchId) {
		branch, _ := vcs.GetBranch(branchId)
		var diffs []vcs.PixelDiff
		if len(branch.Commits) != 0 {
			log.Printf("Old branch!")
			prevImg := composedImage.New(int(branch.Width), int(branch.Height), branch.GetDiffsInBranch())
			diffs = vcs.GetImageDiff(prevImg.Img, *imgRGB)
			log.Printf("Length of pixel diffs %d", len(diffs))
		} else {
			log.Printf("New branch!")
			diffs = vcs.CreateInitialDiff(*imgRGB)

			bounds := imgRGB.Bounds()
			xMax := bounds.Max.X
			xMin := bounds.Min.X
			yMax := bounds.Max.Y
			yMin := bounds.Min.Y

			branch.Width = int16((xMax - xMin))
			branch.Height = int16(yMax - yMin)
		}

		branch.AddCommit(diffs)

		log.Printf("Branch actually expanded? %d, %s", len(branch.Commits), branch.Name)
	} else {
		vcs.CreateOrphanBranch(branchId)
	}
}

func CheckoutCommit(branch string, commit uuid.UUID) composedImage.ComposedImage {

	branchDetails, _ := vcs.GetBranch(branch)
	diffs := make([]vcs.Diff, 0)
	for _, v := range branchDetails.Commits {
		if v == commit {
			break
		}
		diffs = append(diffs, *vcs.GetDiff(v))
	}

	return composedImage.New(int(branchDetails.Width), int(branchDetails.Height), diffs)
}

//func ViewDiff(branchName string, firstCommit uuid.UUID, lastCommit uuid.UUID) []byte {
//	branch := vcs.GetBranch(branchName)
//	commits := branch.GetCommitsRange(firstCommit, lastCommit)
//	pixelDiffs := vcs.SquashCommitsToPixelChanges(commits)
//	// TODO Render pixel diffs and send back
//}

func CreateNewBranch(newBranch string, currentBranch string) {
	_, found := vcs.GetBranch(currentBranch)
	if found == nil {
		log.Println("creating new branch")
		vcs.CreateNewBranch(newBranch, currentBranch)
	} else {
		log.Println("creating orphan branch")
		vcs.CreateOrphanBranch(newBranch)
	}
}

func MergePreview(from string, into string) ([]vcs.Diff, []vcs.Diff) {
	fromBranch, err := vcs.GetBranch(from)
	if err != nil {
		return []vcs.Diff{}, []vcs.Diff{}
	}
	toBranch, err := vcs.GetBranch(into)
	if err != nil {
		return []vcs.Diff{}, []vcs.Diff{}
	}
	toCommits, i, theirDiff, okayDiff, ourDiff := doMerge(fromBranch, toBranch)

	commonCommits := toCommits[:i]

	var diffsTheirs = make([]vcs.Diff, 0)
	var diffsOurs = make([]vcs.Diff, 0)
	for _, v := range commonCommits {
		diffsTheirs = append(diffsTheirs, *vcs.GetDiff(v))
		diffsOurs = append(diffsOurs, *vcs.GetDiff(v))
	}
	diffsTheirs = append(diffsTheirs, vcs.Diff{Commit: uuid.New(), PixelChanges: okayDiff})
	diffsOurs = append(diffsOurs, vcs.Diff{Commit: uuid.New(), PixelChanges: okayDiff})
	diffsTheirs = append(diffsTheirs, vcs.Diff{Commit: uuid.New(), PixelChanges: theirDiff})
	diffsOurs = append(diffsOurs, vcs.Diff{Commit: uuid.New(), PixelChanges: ourDiff})

	return diffsTheirs, diffsOurs
}

func Merge(from string, into string, useTheirs bool) {
	// Find common commit
	fromBranch, err := vcs.GetBranch(from)
	if err != nil {
		return
	}
	toBranch, err := vcs.GetBranch(into)
	if err != nil {
		return
	}
	toCommits, i, theirDiff, okayDiff, ourDiff := doMerge(fromBranch, toBranch)

	toBranch.Commits = toCommits[:i]

	toBranch.AddCommit(okayDiff)
	if useTheirs {
		toBranch.AddCommit(theirDiff)
	} else {
		toBranch.AddCommit(ourDiff)
	}

}

func doMerge(fromBranch *vcs.Branch, toBranch *vcs.Branch) ([]uuid.UUID, int, []vcs.PixelDiff, []vcs.PixelDiff, []vcs.PixelDiff) {
	fromCommits := fromBranch.Commits
	toCommits := make([]uuid.UUID, len(toBranch.Commits))
	copy(toCommits, toBranch.Commits)
	i := 0
	for i < len(fromCommits) && i < len(toCommits) && fromCommits[i] == toCommits[i] {
		i++
	}
	commitsTheirs := fromCommits[i:]
	changesTheirs := vcs.SquashCommitsToPixelChanges(commitsTheirs)
	log.Printf("Theirs has %d new commits with %d ", len(commitsTheirs), len(changesTheirs))
	commitsOurs := toCommits[i:]
	changesOurs := vcs.SquashCommitsToPixelChanges(commitsOurs)
	log.Printf("Ours has %d new commits with %d ", len(commitsOurs), len(changesOurs))
	theirDiff, ourDiff, okayDiff := vcs.AnalyseChanges(changesTheirs, changesOurs)
	log.Printf("Theirs, Ours, Okay pixel changes: %d, %d, %d", len(theirDiff), len(ourDiff), len(okayDiff))
	return toCommits, i, theirDiff, okayDiff, ourDiff
}

func GetBranchNames() []string {
	var branches = vcs.GetBranchHolder()
	var branchNames []string
	fmt.Println(branches.Branches)
	i := 0
	for k := range branches.Branches {
		branchNames = append(branchNames, k)
		i++
	}
	log.Println(branchNames)
	return branchNames
}
