package core

type Position struct {
	LineNo int
	Lat    int // line at 行での位置
	Wat    int // whole? at 全体での位置
}

func (p *Position) Clone() *Position {
	return &Position{
		LineNo: p.LineNo,
		Lat:    p.Lat,
		Wat:    p.Wat,
	}
}

func NewPosition(ln, lat, wat int) *Position {
	return &Position{
		LineNo: ln,
		Lat:    lat,
		Wat:    wat,
	}
}
