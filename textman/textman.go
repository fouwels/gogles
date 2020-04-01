package textman

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png" // Load png decoder
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	gl "github.com/kaelanfouwels/gogles/glow/gl"
)

//Textman Texture Manager
type Textman struct {
	textCache map[string]Texture
}

//Texture ..
type Texture struct {
	Width  int
	Height int
	ID     uint32
}

//NewTextman ..
func NewTextman(textFolder string) (*Textman, error) {

	tman := Textman{
		textCache: map[string]Texture{},
	}

	files, err := ioutil.ReadDir(textFolder)
	if err != nil {
		return nil, fmt.Errorf("Failed to enumerate textFolder %v: %w", textFolder, err)
	}

	for _, f := range files {

		file, err := os.Open(fmt.Sprintf("%v/%v", textFolder, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("Failed to open texture %v: %w", f, err)
		}
		text, err := loadTexture(file)
		if err != nil {
			return nil, fmt.Errorf("Failed to load texture %v: %w", f, err)
		}

		extension := filepath.Ext(f.Name())
		tman.textCache[strings.TrimSuffix(f.Name(), extension)] = text
	}

	return &tman, nil
}

//GetText ...
func (t *Textman) GetText(textName string) (Texture, error) {

	if text, ok := t.textCache[textName]; ok {
		return text, nil
	}
	return Texture{}, fmt.Errorf("Texture %s not present in texture cache", textName)
}

//Destroy ..
func (t *Textman) Destroy() {
	for _, v := range t.textCache {
		gl.DeleteTextures(1, &v.ID)
	}
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
