package renderman

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png" // Load png decoder
	"io"

	gl "github.com/kaelanfouwels/gogles/glow/gl"
)

//Texture ..
type Texture struct {
	Width  int
	Height int
	ID     uint32
}

func loadTexture(reader io.Reader) (Texture, error) {

	img, _, err := image.Decode(reader)
	if err != nil {
		return Texture{}, fmt.Errorf("Failed to decode image: %w", err)
	}

	textRGBA := image.NewRGBA(img.Bounds())
	if textRGBA.Stride != textRGBA.Rect.Size().X*4 {
		return Texture{}, fmt.Errorf("Unsupported stride")
	}
	draw.Draw(textRGBA, textRGBA.Bounds(), img, image.Point{0, 0}, draw.Src)

	size := textRGBA.Rect.Size()

	var tid uint32
	gl.GenTextures(1, &tid)
	gl.BindTexture(gl.TEXTURE_2D, tid)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(textRGBA.Pix))

	return Texture{
		Width:  size.X,
		Height: size.Y,
		ID:     tid,
	}, nil
}
