package fontman

type font struct {
	Height      int
	Description fontDescription
	Metricts    fontMetrics
	Texture     fontTexture
	Chars       []fontChar
	GLTextureID uint32
}

type fontDescription struct {
	Family string
	Style  string
	Size   int
}

type fontTexture struct {
	Name   string
	Width  int
	Height int
}

type fontMetrics struct {
	Ascender  int
	Descender int
	Height    int
}
