package render

import (
	"github.com/disintegration/gift"
	"image"
	"image/color"
	"math"
)

type Modifiable interface {
	Renderable
	FlipX()
	FlipY()
	ApplyColor(c color.Color)
	Copy() Modifiable
	FillMask(img image.RGBA)
	ApplyMask(img image.RGBA)
	Rotate(degrees int)
	Scale(xRatio float64, yRatio float64)
}

func FlipX(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(w-x, y))
		}
	}
	return newRgba
}

func FlipY(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(x, h-y))
		}
	}
	return newRgba
}

func ApplyColor(rgba *image.RGBA, c color.Color) *image.RGBA {
	r1, g1, b1, a1 := c.RGBA()
	r1 = r1 / 257
	g1 = g1 / 257
	b1 = b1 / 257
	a1 = a1 / 257
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r2, g2, b2, a2 := rgba.At(x, y).RGBA()
			r2 = r2 / 257
			g2 = g2 / 257
			b2 = b2 / 257
			a2 = a2 / 257
			a3 := a1 + a2
			tmp := color.RGBA{
				uint8(((a1 * r1) + (a2 * r2)) / a3),
				uint8(((a1 * g1) + (a2 * g2)) / a3),
				uint8(((a1 * b1) + (a2 * b2)) / a3),
				uint8(a2)}
			newRgba.Set(x, y, tmp)
		}
	}
	return newRgba
}

func FillMask(rgba *image.RGBA, img image.RGBA) *image.RGBA {
	// Instead of static color it just two buffers melding
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(x, y))
		}
	}
	bounds = img.Bounds()
	w = bounds.Max.X
	h = bounds.Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r1, g1, b1, a1 := rgba.At(x, y).RGBA()
			r2, g2, b2, a2 := img.At(x, y).RGBA()

			var tmp color.RGBA

			a1 = a1 / 257
			a2 = a2 / 257

			r1 = r1 / 257
			g1 = g1 / 257
			b1 = b1 / 257

			r2 = r2 / 257
			g2 = g2 / 257
			b2 = b2 / 257

			if a1 == 0 {
				tmp = color.RGBA{
					uint8(r2),
					uint8(g2),
					uint8(b2),
					uint8(a2),
				}
			} else {
				tmp = color.RGBA{
					uint8(r1),
					uint8(g1),
					uint8(b1),
					uint8(a1),
				}
			}

			newRgba.Set(x, y, tmp)
		}
	}
	return newRgba
}

func ApplyMask(rgba *image.RGBA, img image.RGBA) *image.RGBA {
	// Instead of static color it just two buffers melding
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(x, y))
		}
	}
	bounds = img.Bounds()
	w = bounds.Max.X
	h = bounds.Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r1, g1, b1, a1 := img.At(x, y).RGBA()
			r2, g2, b2, a2 := rgba.At(x, y).RGBA()

			var tmp color.RGBA

			a1 = a1 / 257
			a2 = a2 / 257
			a3 := a1 + a2
			if a3 == 0 {
				tmp = color.RGBA{
					0, 0, 0, 0,
				}
				newRgba.Set(x, y, tmp)
				continue
			}

			r1 = r1 / 257
			g1 = g1 / 257
			b1 = b1 / 257

			r2 = r2 / 257
			g2 = g2 / 257
			b2 = b2 / 257

			tmp = color.RGBA{
				uint8(((a1 * r1) + (a2 * r2)) / a3),
				uint8(((a1 * g1) + (a2 * g2)) / a3),
				uint8(((a1 * b1) + (a2 * b2)) / a3),
				uint8(math.Max(float64(a1), float64(a2)))}

			newRgba.Set(x, y, tmp)
		}
	}
	return newRgba
}

func Rotate(rgba *image.RGBA, degrees int) *image.RGBA {
	filter := gift.New(
		gift.Rotate(float32(degrees), color.Black, gift.CubicInterpolation))
	dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
	filter.Draw(dst, rgba)
	return dst

}

func Scale(rgba *image.RGBA, xRatio float64, yRatio float64) *image.RGBA {
	bounds := rgba.Bounds()
	w := int(math.Floor(float64(bounds.Max.X) * xRatio))
	h := int(math.Floor(float64(bounds.Max.Y) * yRatio))
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(int(math.Floor(float64(x)/xRatio)), int(math.Floor(float64(y)/yRatio))))
		}
	}
	return newRgba
}

func round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}