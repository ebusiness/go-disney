package utils

import (
	"github.com/llgcode/draw2d"
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
	"wait20":     "verysatisfied.png",
	"wait30":     "satisfied.png",
	"wait40":     "neutral.png",
	"wait50":     "dissatisfied.png",
	"wait60":     "verydissatisfied.png",
	"showwait20": "showverysatisfied.png",
	"showwait30": "showsatisfied.png",
	"showwait40": "showneutral.png",
	"showwait50": "showdissatisfied.png",
	"showwait60": "showverydissatisfied.png",
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

// DrawString -
func (pd *PlanDraw) DrawString(str string, point DrawPoint) *PlanDraw {
	gc := draw2dimg.NewGraphicContext(pd.dst)
	draw2d.SetFontFolder("/asset/font")
	gc.SetFontData(draw2d.FontData{Name: "luxi", Style: draw2d.FontStyleBold})
	gc.SetFillColor(image.White)
	// gc.SetFillColor(image.NewUniform(color.RGBA{0x21, 0x96, 0xf3, 0xff}))
	gc.SetFontSize(24)
	gc.FillStringAt(str, point.X, point.Y)
	// DrawString(pd.dst, int(point.X), int(point.Y), str, color.White)
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
	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0x88})
	gc.SetLineWidth(5)

	for index, line := range lines {
		if index == 0 {
			gc.MoveTo(line.X, line.Y) // should always be called first for a new path
			continue
		}
		gc.LineTo(line.X, line.Y)
	}
	// gc.Close()
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
