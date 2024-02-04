package backend

import (
	"collaborart/backend/vcs"
	"os"
)

func PushToBranch(branch string, image *os.File) {
	if vcs.BranchExists(branch) {
		//var changes []vcs.PixelDiff = // TODO Get diff between this image and tip of branch

	} else {
		vcs.CreateOrphanBranch(branch)
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
