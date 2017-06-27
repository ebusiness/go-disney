package utils

import (
	"github.com/llgcode/draw2d/draw2dimg"
	log "github.com/sirupsen/logrus"
	"image"
	"image/color"
	"image/draw"
	"image/png" // register the PNG format with the image package
	// "log"
	// "math"
	"os"
)

var planResources = map[string]string{
	"land":       "disneyland.png",
	"sea":        "disneysea.png",
	"wait20":     "blue.png",
	"wait30":     "green.png",
	"wait40":     "orange.png",
	"wait50":     "deeporange.png",
	"wait60":     "red.png",
	"wait70":     "deepred.png",
	"showwait20": "show_blue.png",
	"showwait30": "show_green.png",
	"showwait40": "show_orange.png",
	"showwait50": "show_deeporange.png",
	"showwait60": "show_red.png",
	"showwait70": "show_deepred.png",
}

// load an image
func loadImage(resID string) image.Image {
	res := planResources[resID]
	if len(res) < 1 {
		log.Fatalln("resource not found")
	}
	path := "/asset/base/" + res
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Fatalf("can not load file [%s]", path)
	}
	img, err := png.Decode(file)

	if err != nil {
		log.Fatalf("[%s] is not a image ", path)
	}

	return img
}

// NewPlanDraw make a draw with bacground image name
func NewPlanDraw(resID string) *PlanDraw {
	img := loadImage(resID)

	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)

	return &PlanDraw{m}
}

// DrawPoint -
type DrawPoint struct {
	X, Y float64
}

// PlanDraw -
type PlanDraw struct {
	dst draw.Image
}

// DrawMark -
func (pd *PlanDraw) DrawMark(resID string, point DrawPoint) *PlanDraw {
	mark := loadImage(resID)
	offset := image.Pt(int(point.X)-mark.Bounds().Max.X/2, int(point.Y)-mark.Bounds().Max.Y)
	draw.Draw(pd.dst, mark.Bounds().Add(offset), mark, image.ZP, draw.Over)

	return pd
}

// DrawLines -
func (pd *PlanDraw) DrawLines(lines []DrawPoint) *PlanDraw {
	if len(lines) == 0 {
		return pd
	}
	gc := draw2dimg.NewGraphicContext(pd.dst)
	// Set some properties
	// gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.Black)
	gc.SetLineWidth(5)

	for index, line := range lines {
		if index == 0 {
			gc.MoveTo(line.X, line.Y) // should always be called first for a new path
			continue
		}
		gc.LineTo(line.X, line.Y)
	}
	gc.Close()
	gc.Stroke()

	return pd
}

// SaveImage -
func (pd *PlanDraw) SaveImage(filepath string) {
	imgw, err := os.Create(filepath)
	defer imgw.Close()

	if err != nil {
		log.Fatalf("cannot save [%s]", filepath)
	}
	png.Encode(imgw, pd.dst)
}
