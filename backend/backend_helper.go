package backend

import (
	"collaborart/backend/composedImage"
	"collaborart/backend/vcs"
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
			prevImg := composedImage.New(branch)
			diffs = vcs.GetImageDiff(prevImg.Img, *imgRGB)
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

//func CheckoutCommit(branch string, commit uuid.UUID) []byte {
//
//}

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

//func Merge(from string, into string) {
//	// Find common commit
//	fromBranch, err := vcs.GetBranch(from)
//	if err != nil {
//		return
//	}
//	toBranch, err := vcs.GetBranch(into)
//	if err != nil {
//		return
//	}
//	fromCommits := fromBranch.Commits
//	toCommits := toBranch.Commits
//	i := 0
//	for i < len(fromCommits) && i < len(toCommits) && fromCommits[i] == toCommits[i] {
//		i++
//	}
//	commitsTheirs := fromCommits[i:]
//	changesTheirs := vcs.SquashCommitsToPixelChanges(commitsTheirs)
//	commitsOurs := toCommits[i:]
//	changesOurs := vcs.SquashCommitsToPixelChanges(commitsOurs)
//	theirDiff, ourDiff, okayDiff := vcs.AnalyseChanges(changesTheirs, changesOurs)
//
//}

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
