package gopdf

import (
	"fmt"
	"io"
)

type cacheContentImage struct {
	verticalFlip   bool
	horizontalFlip bool
	index          int
	x              float64
	y              float64
	h              float64
	rect           Rect
	transparency   *Transparency
	crop           *CropOptions
}

func (c *cacheContentImage) write(w io.Writer, protection *PDFProtection) error {
	x := c.x
	width := c.rect.W
	height := c.rect.H
	y := c.h - (c.y + c.rect.H)

	contentStream := "q\n"

	if c.transparency != nil && c.transparency.Alpha != 1 {
		contentStream += fmt.Sprintf("/GS%d gs\n", c.transparency.indexOfExtGState)
	}

	if c.horizontalFlip || c.verticalFlip {
		fh := "1"
		if c.horizontalFlip {
			fh = "-1"
			x = -1*x - width
		}

		fv := "1"
		if c.verticalFlip {
			fv = "-1"
			y = -1*y - height
		}

		contentStream += fmt.Sprintf("%s 0 0 %s 0 0 cm\n", fh, fv)
	}

	if c.crop != nil {
		contentStream += fmt.Sprintf("\t\t%0.2f %0.2f %0.2f %0.2f re W* n\n", x, y, c.crop.Width, c.crop.Height)
		contentStream += fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", width, height, x - c.crop.X, y - c.crop.Y, c.index+1)
	} else {
		contentStream += fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", width, height, x, y, c.index+1)
	}

	contentStream += "Q\n"

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}
