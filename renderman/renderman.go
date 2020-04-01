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
	gl.ClearColor(0, 0, 0, 0)
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
func (r *RenderMan) Draw() error {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	err := r.drawBackground()
	if err != nil {
		return err
	}
	err = r.drawForeground()
	if err != nil {
		return err
	}
	err = r.drawOverlay()
	if err != nil {
		return err
	}

	return nil
}

func (r *RenderMan) drawBackground() error {

	text, err := r.textman.GetText("t_square")
	if err != nil {
		return err
	}

	gl.Enable(gl.TEXTURE_2D)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.REPLACE)
	gl.BindTexture(gl.TEXTURE_2D, text.ID)

	gl.MatrixMode(gl.TEXTURE)
	gl.LoadIdentity()

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(-200, -50, 0)
	gl.Rotatef(45, 0, 0, 1)
	gl.Scalef(100, 100, 0)

	gl.Begin(gl.QUADS)
	gl.Color3f(1.0, 1.0, 1.0)

	gl.TexCoord2f(0, 0)
	gl.Vertex2f(0, 0)

	gl.TexCoord2f(0, 1)
	gl.Vertex2f(0, 1)

	gl.TexCoord2f(1, 1)
	gl.Vertex2f(1, 1)

	gl.TexCoord2f(1, 0)
	gl.Vertex2f(1, 0)

	gl.End()

	return nil
}

func (r *RenderMan) drawForeground() error {
	return nil
}

func (r *RenderMan) drawOverlay() error {

	err := r.fontman.RenderChar('M', 0, 0, 1)
	if err != nil {
		return err
	}
	err = r.fontman.RenderString("Fouwels", 0, -100, 0.25)
	if err != nil {
		return err
	}

	return nil
}
