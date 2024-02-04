package backend

import (
	"collaborart/backend/composedImage"
	"collaborart/backend/vcs"
	"image"
	"image/jpeg"
	"log"
)

func PushToBranch(branchId string, imageFile *image.Image) {
	img, _ := jpeg.Decode(imageFile)
	var imgRGB image.RGBA
	if x, ok := img.(*image.RGBA); ok {
		imgRGB = *x
	}

	if vcs.BranchExists(branchId) {
		branch, _ := vcs.GetBranch(branchId)
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
	_, found := vcs.GetBranch(currentBranch)
	if found == nil {
		log.Println("creating new branch")
		vcs.CreateNewBranch(newBranch, currentBranch)
	} else {
		log.Println("creating orphan branch")
		vcs.CreateOrphanBranch(newBranch)
	}
}

func Merge(from string, into string) {
	// Find common commit
	fromBranch, err := vcs.GetBranch(from)
	if err != nil {
		return
	}
	toBranch, err := vcs.GetBranch(into)
	if err != nil {
		return
	}
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
