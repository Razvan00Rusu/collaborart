package backend

import (
	"collaborart/backend/composedImage"
	"collaborart/backend/vcs"
	"image"
	"image/jpeg"
	"os"
)

func PushToBranch(branchId string, imageFile *os.File) {
	img, _ := jpeg.Decode(imageFile)
	var imgRGB image.RGBA
	if x, ok := img.(*image.RGBA); ok {
		imgRGB = *x
	}

	if vcs.BranchExists(branchId) {
		branch := vcs.GetBranch(branchId)
		var diffs []vcs.PixelDiff
		if len(branch.Commits) != 0 {
			prevImg := composedImage.New(branch)
			diffs = vcs.GetImageDiff(prevImg.Img, imgRGB)
		} else {
			diffs = vcs.CreateInitialDiff(imgRGB)
		}

		branch.AddCommit(diffs)
	} else {
		vcs.CreateOrphanBranch(branchId)
	}
}

//func CheckoutCommit(branch string, commit uuid.UUID) []byte {
//
//}

func CreateNewBranch(newBranch string, currentBranch string) {
	vcs.CreateNewBranch(newBranch, currentBranch)
}

func Merge(from string, into string) {
	// Find common commit
	fromBranch := vcs.GetBranch(from)
	toBranch := vcs.GetBranch(into)
	fromCommits := fromBranch.Commits
	toCommits := toBranch.Commits
	i := 0
	for i < len(fromCommits) && i < len(toCommits) && fromCommits[i] == toCommits[i] {
		i++
	}
	commitsToAdd := fromCommits[i:]
	changes := vcs.SquashCommitsToPixelChanges(commitsToAdd)
	toBranch.AddCommit(changes)
}

func GetBranchNames() []string {
	var branches = vcs.GetBranchHolder()
	var branchNames = make([]string, len(branches.Branches))
	i := 0
	for k := range branches.Branches {
		branchNames[i] = k
		i++
	}
	return branchNames
}
