package renderman

import (
	"fmt"
	"os"

	"github.com/kaelanfouwels/gogles/fontman"
	gl "github.com/kaelanfouwels/gogles/glow/gl"
)

const assetsDir string = "assets/"

//RenderMan ..
type RenderMan struct {
	fontman   *fontman.Fontman
	width     float32
	height    float32
	textCache map[string]Texture
}

//NewRenderman ..
func NewRenderman(width float32, height float32, fontman *fontman.Fontman) (*RenderMan, error) {

	rm := RenderMan{
		width:     width,
		height:    height,
		fontman:   fontman,
		textCache: map[string]Texture{},
	}

	rm.initialize()

	textures := []string{
		"t_square",
		"t_background",
	}

	for _, v := range textures {
		f, err := os.Open(fmt.Sprintf("%v/%v.png", assetsDir, v))
		if err != nil {
			return nil, fmt.Errorf("Failed to open %v: %w", v, err)
		}

		text, err := loadTexture(f)
		if err != nil {
			return nil, fmt.Errorf("Failed to load %v: %w", v, err)
		}
		rm.textCache[v] = text
	}

	return &rm, nil
}

//Destroy ..
func (r *RenderMan) Destroy() {
	for _, v := range r.textCache {
		gl.DeleteTextures(1, &v.ID)
	}
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
	//r.drawBackground()
	r.drawForeground()
	r.drawOverlay()
}

func (r *RenderMan) drawBackground() {

	gl.Enable(gl.TEXTURE_2D)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.REPLACE)
	gl.BindTexture(gl.TEXTURE_2D, r.textCache["t_square"].ID)

	gl.Begin(gl.QUADS)
	gl.Color3f(1.0, 1.0, 1.0)

	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-r.width/2, -r.height/2)

	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-r.width/2, r.height/2)

	gl.TexCoord2f(1, 1)
	gl.Vertex2f(r.width/2, r.height/2)

	gl.TexCoord2f(1, 0)
	gl.Vertex2f(r.width/2, -r.height/2)

	gl.End()
}

func (r *RenderMan) drawForeground() {

}

func (r *RenderMan) drawOverlay() {
	r.fontman.RenderChar('F', 0, 0)
}
