package fontman

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png" // Load png decoder
	"os"

	gl "github.com/kaelanfouwels/gogles/glow/gl"
)

//Fontman Font Manager
type Fontman struct {
	textID uint32
	font   font
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
func NewFontman() (*Fontman, error) {
	fm := Fontman{
		font: consolasRegular24,
	}

	filename := fm.font.File

	ffile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to load font file %s: %w", filename, err)
	}
	img, _, err := image.Decode(ffile)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode font file %s: %w", filename, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride in %s", filename)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var textID uint32
	gl.GenTextures(1, &textID)
	gl.BindTexture(gl.TEXTURE_2D, textID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.LINEAR)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	fm.textID = textID
	return &fm, nil
}

//RenderChar Render a character
func (f *Fontman) RenderChar(char rune, x float32, y float32) error {

	fchar, err := f.font.LookupFontChar(char)
	if err != nil {
		return err
	}

	gl.Enable(gl.TEXTURE_2D)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
	gl.BindTexture(gl.TEXTURE_2D, f.textID)

	gl.Begin(gl.QUADS)
	gl.Color4f(1.0, 1.0, 1.0, 1.0)

	cheight := fchar.Y - fchar.H
	cwidth := fchar.W - fchar.X

	//0,0
	gl.TexCoord2f(0, fchar.Y)
	gl.Vertex2f(x, y)

	//0,1
	gl.TexCoord2f(fchar.X, fchar.H)
	gl.Vertex2f(x, y+cheight)

	//1,1
	gl.TexCoord2f(fchar.W, fchar.H)
	gl.Vertex2f(x+cwidth, y+cheight)

	//1,0
	gl.TexCoord2f(fchar.W, fchar.Y)
	gl.Vertex2f(x+cwidth, y+cheight)

	gl.End()

	return nil
}
