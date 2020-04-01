package fontman

import (
	"fmt"
	_ "image/png" // Load png decoder

	gl "github.com/kaelanfouwels/gogles/glow/gl"
	"github.com/kaelanfouwels/gogles/textman"
)

//Fontman Font Manager
type Fontman struct {
	textman *textman.Textman
	font    font
}

func (f *font) LookupFontChar(char rune) (fontChar, error) {
	// No this is not efficient
	for _, c := range f.Chars {
		cl := c
		if cl.Char == char {
			return cl, nil
		}
	}

	return fontChar{}, fmt.Errorf("Char %v not found in font index", char)
}

type fontChar struct {
	Char  rune
	Width int
	X     float32
	Y     float32
	W     float32
	H     float32
	OX    float32
	OY    float32
}

//NewFontman Generate a new font manager
func NewFontman(textman *textman.Textman) (*Fontman, error) {
	fm := Fontman{
		font:    consolasRegular24,
		textman: textman,
	}

	return &fm, nil
}

//RenderChar Render a character
func (f *Fontman) RenderChar(char rune, x float32, y float32) error {

	fchar, err := f.font.LookupFontChar(char)
	if err != nil {
		return err
	}

	_ = fchar

	ftext, err := f.textman.GetText(f.font.Texture.Name)
	if err != nil {
		return err
	}

	gl.Enable(gl.TEXTURE_2D)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.REPLACE)
	gl.BindTexture(gl.TEXTURE_2D, ftext.ID)

	gl.Begin(gl.QUADS)
	gl.Color3f(1.0, 1.0, 1.0)

	//0,0
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-256/2, -256/2)

	//0,1
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-256/2, 256/2)

	//1,1
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(256/2, 256/2)

	//1,0
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(256/2, -256/2)

	gl.End()
	gl.PopMatrix()

	return nil
}
