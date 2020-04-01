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
		font:    consolasRegular65,
		textman: textman,
	}

	return &fm, nil
}

//RenderString ..
func (f *Fontman) RenderString(text string, x float32, y float32, scaling float32) error {

	const kerning float32 = 10

	rs := []rune(text)

	xCursor := x

	for _, v := range rs {

		err := f.RenderChar(v, xCursor, y, scaling)
		if err != nil {
			return err
		}

		fchar, err := f.font.LookupFontChar(v)
		if err != nil {
			return err
		}

		xCursor += float32(fchar.W) + kerning
	}

	return nil
}

//RenderChar Render a character
func (f *Fontman) RenderChar(char rune, x float32, y float32, scaling float32) error {

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
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
	gl.BindTexture(gl.TEXTURE_2D, ftext.ID)

	gl.MatrixMode(gl.TEXTURE)
	gl.LoadIdentity()
	gl.Scalef(1/float32(ftext.Width), 1/float32(ftext.Height), 1)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Scalef(scaling, scaling, 0)

	gl.Begin(gl.QUADS)

	woff := fchar.W
	hoff := fchar.H
	//0,0
	gl.TexCoord2f(fchar.X, fchar.Y)
	gl.Vertex2f(x, y+hoff)

	//0,1
	gl.TexCoord2f(fchar.X, fchar.Y+hoff)
	gl.Vertex2f(x, y)

	//1,1
	gl.TexCoord2f(fchar.X+woff, fchar.Y+hoff)
	gl.Vertex2f(x+woff, y)

	//1,0
	gl.TexCoord2f(fchar.X+woff, fchar.Y)
	gl.Vertex2f(x+woff, y+hoff)

	gl.End()

	return nil
}
