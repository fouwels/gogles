package fontman

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	gl "github.com/kaelanfouwels/gogles/glow/gl"
)

//Fontman Font Manager
type Fontman struct {
	textID uint32
	font   font
}

type font struct {
	File        string
	Height      int
	Description struct {
		Family string
		Style  string
		Size   string
	}
	Metricts struct {
		Ascender  int
		Descender int
		Height    int
	}
	Texture struct {
		File   string
		Width  int
		Height int
	}
	Chars []fontChar
}

func (f *font) LookupFontChar(char rune) (fontChar, error) {
	// No this is not efficient
	for _, c := range f.Chars {
		cl := c
		if cl.Char == char {
			return cl, nil
		}
	}

	return fontChar{}, fmt.Errorf("Char %s not found in font index", char)
}

type fontChar struct {
	Char  rune
	Width int
	X     int
	Y     int
	W     int
	H     int
	OX    int
	OY    int
}

//NewFontman Generate a new font manager
func NewFontman() (*Fontman, error) {
	fm := Fontman{
		font: consolas_regular_24,
	}

	ffile, err := os.Open(fmt.Sprintf(fm.font.File))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(ffile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var textID uint32
	gl.Enable(gl.TEXTURE_2D)
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
func (f *Fontman) RenderChar(char rune, centerX int, centerY int) error {

	fchar, err := f.font.LookupFontChar(char)
	if err != nil {
		return err
	}

	gl.Begin(gl.QUADS)

	gl.Vertex3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1)

	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1)

	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)

	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)

	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)

	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)

	gl.Text

	gl.End()
}
