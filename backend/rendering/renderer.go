package rendering

import (
	"collaborart/backend/vcs"
	"image"
	"io"
)

type renderer interface {
	New(w io.Writer) renderer
	compose(branch vcs.Branch) image.Image
	render(img image.Image) error
}
