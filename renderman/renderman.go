package renderman

import (
	"github.com/kaelanfouwels/gogles/fontman"
	gl "github.com/kaelanfouwels/gogles/glow/gl"
	"github.com/kaelanfouwels/gogles/textman"
)

const assetsDir string = "assets/"

//RenderMan ..
type RenderMan struct {
	textman *textman.Textman
	fontman *fontman.Fontman
	width   float32
	height  float32
}

//NewRenderman ..
func NewRenderman(width float32, height float32, textman *textman.Textman, fontman *fontman.Fontman) (*RenderMan, error) {

	rm := RenderMan{
		width:   width,
		height:  height,
		textman: textman,
		fontman: fontman,
	}

	rm.initialize()

	return &rm, nil
}

//Destroy ..
func (r *RenderMan) Destroy() {

}

func (r *RenderMan) initialize() {
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ShadeModel(gl.FLAT)
	gl.Enable(gl.DEPTH_TEST)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-float64(r.width)/2, float64(r.width)/2, -float64(r.height)/2, float64(r.height)/2, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

// Draw ..
func (r *RenderMan) Draw() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	r.drawBackground()
	r.drawForeground()
	r.drawOverlay()
}

func (r *RenderMan) drawBackground() {

	// text, err := r.textman.GetText("t_square")
	// if err != nil {
	// 	panic(err)
	// }

	// gl.Enable(gl.TEXTURE_2D)
	// gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.REPLACE)
	// gl.BindTexture(gl.TEXTURE_2D, text.ID)

	// gl.Begin(gl.QUADS)
	// gl.Color3f(1.0, 1.0, 1.0)

	// gl.TexCoord2f(0, 0)
	// gl.Vertex2f(-r.width/2, -r.height/2)

	// gl.TexCoord2f(0, 1)
	// gl.Vertex2f(-r.width/2, r.height/2)

	// gl.TexCoord2f(1, 1)
	// gl.Vertex2f(r.width/2, r.height/2)

	// gl.TexCoord2f(1, 0)
	// gl.Vertex2f(r.width/2, -r.height/2)

	// gl.End()
}

func (r *RenderMan) drawForeground() {

}

func (r *RenderMan) drawOverlay() {
	r.fontman.RenderChar('F', 0, 0)
}
