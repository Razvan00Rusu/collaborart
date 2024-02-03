package composedImage

import (
	"collaborart/backend/vcs"
	"image"
	"image/color"
)

type composedImage struct {
	img image.RGBA
}

func New(branch vcs.Branch) composedImage {

	// TODO: build the image with size from the branch
	picture := image.NewRGBA(image.Rect(0, 0, 8, 5))

	for _, change := range branch.GetDiffsInBranch() {

		for _, diff := range change.PixelChanges {
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

	return composedImage{img: *picture}
}
