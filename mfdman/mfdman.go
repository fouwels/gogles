package mfdman

import (
	"github.com/kaelanfouwels/gogles/fontman"
	"github.com/kaelanfouwels/gogles/glow/gl"
)

//MFDIndex defines an MFD index
type MFDIndex int

const mfdWidth float32 = 80
const mfdHeight float32 = 80
const mfdYOffset float32 = 60
const mfdXOffset float32 = 20

const (

	// L1 ..
	L1 MFDIndex = iota
	// L2 ..
	L2
	// L3 ..
	L3
	// L4 ..
	L4
	// R1 ..
	R1
	// R2 ..
	R2
	// R3 ..
	R3
	// R4 ..
	R4
	// MFDCount ..
	MFDCount
)

type mfd struct {
	x        float32
	y        float32
	selected bool
	textA    string
	textB    string
}

//MFDman ..
type MFDman struct {
	width   float32
	height  float32
	mfds    [MFDCount]mfd
	fontman *fontman.Fontman
}

//NewMFDman ..
func NewMFDman(width float32, height float32, fontman *fontman.Fontman) (*MFDman, error) {

	mfdm := MFDman{
		width:   width,
		height:  height,
		fontman: fontman,
	}

	ycursor := -height/2 + mfdYOffset
	for i := 0; i < int(MFDCount/2); i++ {
		mfdm.mfds[i].y = ycursor
		mfdm.mfds[i].x = -width/2 + mfdXOffset
		mfdm.mfds[i].textA = ""
		mfdm.mfds[i].textB = ""
		mfdm.mfds[i].selected = false
		ycursor += mfdHeight
		ycursor += mfdYOffset
	}

	ycursor = -height/2 + mfdYOffset
	for i := int(MFDCount / 2); i < int(MFDCount); i++ {
		mfdm.mfds[i].y = ycursor
		mfdm.mfds[i].x = width/2 - mfdWidth - mfdXOffset
		mfdm.mfds[i].textA = ""
		mfdm.mfds[i].textB = ""
		mfdm.mfds[i].selected = false
		ycursor += mfdHeight
		ycursor += mfdYOffset
	}

	return &mfdm, nil
}

//Draw ..
func (m *MFDman) Draw() error {

	for _, v := range m.mfds {
		err := m.drawOne(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MFDman) drawOne(mfd mfd) error {

	if mfd.textA == "" && mfd.textB == "" {
		return nil
	}
	gl.MatrixMode(gl.TEXTURE)
	gl.LoadIdentity()
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Disable(gl.TEXTURE_2D)

	// Draw MFD box
	gl.LineWidth(3)
	if !mfd.selected {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	gl.Begin(gl.QUADS)

	//0,0
	gl.Vertex2f(mfd.x, mfd.y)

	//0,1
	gl.Vertex2f(mfd.x, mfd.y+mfdHeight)

	//1,1
	gl.Vertex2f(mfd.x+mfdWidth, mfd.y+mfdHeight)

	//1,0
	gl.Vertex2f(mfd.x+mfdWidth, mfd.y)
	gl.End()

	// Draw MFD legend
	ycursor := mfd.y + mfdHeight - 10
	ycursor -= 20
	err := m.fontman.RenderString(mfd.textA, mfd.x+10, ycursor, 0.20)
	if err != nil {
		return err
	}
	ycursor -= 15
	ycursor -= 20
	err = m.fontman.RenderString(mfd.textB, mfd.x+10, ycursor, 0.20)
	if err != nil {
		return err
	}

	return nil
}

//SetText ..
func (m *MFDman) SetText(mfd MFDIndex, textA string, textB string) {
	m.mfds[mfd].textA = textA
	m.mfds[mfd].textB = textB
}

//SetSelected ..
func (m *MFDman) SetSelected(mfd MFDIndex, selected bool) {
	m.mfds[mfd].selected = selected
}

//GetSelected ..
func (m *MFDman) GetSelected(mfd MFDIndex) bool {
	return m.mfds[mfd].selected
}
