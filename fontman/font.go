package fontman

type font struct {
	File        string
	Height      int
	Description fontDescription
	Metricts    fontMetrics
	Texture     fontTexture
	Chars       []fontChar
}

type fontDescription struct {
	Family string
	Style  string
	Size   int
}

type fontTexture struct {
	File   string
	Width  int
	Height int
}

type fontMetrics struct {
	Ascender  int
	Descender int
	Height    int
}
