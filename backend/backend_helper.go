package backend

import (
	"collaborart/backend/vcs"
)

//func PushToBranch(branch string, image []byte) {
//	if vcs.BranchExists(branch) {
//		var changes []vcs.PixelDiff = // TODO Get diff between this image and tip of branch
//
//	} else {
//		vcs.CreateOrphanBranch(branch)
//	}
//}

//func CheckoutCommit(branch string, commit uuid.UUID) []byte {
//
//}

func CreateNewBranch(newBranch string, currentBranch string) {
	vcs.CreateNewBranch(newBranch, currentBranch)
}

func Merge(from string, into string) {

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
