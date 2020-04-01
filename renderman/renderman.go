package renderman

import (
	"github.com/kaelanfouwels/gogles/fontman"
	gl "github.com/kaelanfouwels/gogles/glow/gl"
)

//RenderMan ..
type RenderMan struct {
	fontMan *fontman.FontMan
	width   int
	height  int
}

//NewRenderman ..
func NewRenderman(width int, height int, fontman *fontman.Fontman) *RenderMan {

	rm := &RenderMan{
		width:   width,
		height:  height,
		fontman: fontman,
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.5, 0.5, 0.5, 1}
	diffuse := []float32{1, 1, 1, 1}
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	f := ((float64(width) / float64(height)) - 1) / 2
	gl.Frustum(-1-f, 1+f, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	return
}

//Destroy ..
func (r *RenderMan) Destroy() {

}

// Draw ..
func (r *RenderMan) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.fontMan.RenderChar('F', r.width/2, r.height/2)
}
