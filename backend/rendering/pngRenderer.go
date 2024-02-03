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

		for _, diff := range change {
			picture.Set(5, 5, color.RGBA{255, 0, 0, 255})
		}

	}

	return *picture
}
