package rendering

import (
	"collaborart/backend/vcs"
	"image"
	"image/color"
	"io"
)

type PngRenderer struct {
	w io.Writer
}

func New(w io.Writer) PngRenderer {
	return PngRenderer{w}
}

func compose(branch vcs.Branch) image.RGBA {
	// TODO: build the image with size from the branch
	picture := image.NewRGBA(image.Rect(0, 0, 8, 5))

	for _, change := range branch.GetDiffsInBranch() {

		for _, diff := range change.PixelChanges {
			// I think there will be an issue with Alpha, but just ignore
			x := int(diff.X)
			y := int(diff.Y)
			oldR, oldG, oldB, oldA := picture.At(x, y).RGBA()
			newColor := color.RGBA{
				uint8(int16(oldR) + diff.DR),
				uint8(int16(oldG) + diff.DG),
				uint8(int16(oldB) + diff.DB),
				uint8(int16(oldA) + diff.DA),
			}

			picture.Set(
				x,
				y,
				newColor,
			)
		}
	}

	// kinda cringe to copy this, change if relevant
	return *picture
}
