package composedImage

import (
	"collaborart/backend/vcs"
	"image"
	"image/color"
	"log"
)

type composedImage struct {
	Img image.RGBA
}

func New(branch *vcs.Branch) composedImage {

	// TODO: build the image with size from the branch
	picture := image.NewRGBA(image.Rect(0, 0, int(branch.Width), int(branch.Height)))

	branchDiffs := branch.GetDiffsInBranch()
	log.Printf("Image combosed from branch? %d, %s, %d, %d", len(branch.Commits), branch.Name, len(branchDiffs), len(branchDiffs[0].PixelChanges))

	for _, change := range branch.GetDiffsInBranch() {

		//log.Printf("A change")

		for _, diff := range change.PixelChanges {

			//log.Printf("A diff")

			// I think there will be an issue with Alpha, but just ignore
			x := int(diff.X)
			y := int(diff.Y)
			newColor := color.RGBA{
				R: uint8(diff.R),
				G: uint8(diff.G),
				B: uint8(diff.B),
				A: uint8(diff.A),
			}

			picture.Set(
				x,
				y,
				newColor,
			)
		}
	}

	return composedImage{Img: *picture}
}
