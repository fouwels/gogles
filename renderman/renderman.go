package renderman

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/kaelanfouwels/gogles/ioman"
	"github.com/kaelanfouwels/gogles/mfdman"
)

const assetsDir string = "assets/"

//RenderMan ..
type RenderMan struct {
	mfdman *mfdman.MFDman
	ioman  *ioman.IOMan
	width  float32
	height float32
	gc     *gg.Context
}

//NewRenderman ..
func NewRenderman(gc *gg.Context, width float32, height float32, mfdman *mfdman.MFDman, ioman *ioman.IOMan) (*RenderMan, error) {

	rm := RenderMan{
		gc:     gc,
		width:  width,
		height: height,
		mfdman: mfdman,
		ioman:  ioman,
	}

	rm.initialize()

	return &rm, nil
}

//Destroy ..
func (r *RenderMan) Destroy() {

}

func (r *RenderMan) initialize() {

}

// Draw ..
func (r *RenderMan) Draw() error {

	err := r.drawBackground()
	if err != nil {
		return err
	}
	err = r.drawForeground()
	if err != nil {
		return err
	}
	err = r.mfdman.Draw()
	if err != nil {
		return err
	}

	return nil
}

var testCounter float32 = 0

func (r *RenderMan) drawBackground() error {

	r.gc.SetColor(color.RGBA{25, 25, 25, 255})
	r.gc.Clear()

	testCounter++
	return nil
}

func (r *RenderMan) drawForeground() error {

	return nil
}
