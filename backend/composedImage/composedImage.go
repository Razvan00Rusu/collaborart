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

	for _, change := range branch.GetDiffsInBranch() {

		log.Printf("A change")

		for _, diff := range change.PixelChanges {

			log.Printf("A diff")

			// I think there will be an issue with Alpha, but just ignore
			x := int(diff.X)
			y := int(diff.Y)
			oldR, oldG, oldB, oldA := picture.At(x, y).RGBA()
			newColor := color.RGBA{
				R: uint8(int16(oldR) + diff.DR),
				G: uint8(int16(oldG) + diff.DG),
				B: uint8(int16(oldB) + diff.DB),
				A: uint8(int16(oldA) + diff.DA),
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
