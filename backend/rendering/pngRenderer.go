package rendering

import (
	"collaborart/backend/vcs"
	"image"
	"io"
)

type PngRenderer struct {
	w io.Writer
}

func New(w io.Writer) PngRenderer {
	return PngRenderer{w}
}

func compose(branch vcs.Branch) {
	// TODO: build the image with size from the branch
	picture := image.NewRGBA(image.Rect(0, 0, 8, 5))

}
