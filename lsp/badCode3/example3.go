package badCode3

type ConvexQuadrilateral interface {
	GetArea() int
}

type Rectangle interface {
	ConvexQuadrilateral
	SetA(a int)
	SetB(b int)
}

type Oblong struct {
	Rectangle
	a int
	b int
}

func (o *Oblong) SetA(a int) {
	o.a = a
}

func (o *Oblong) SetB(b int) {
	o.b = b
}

type Square struct {
	Rectangle
	a int
}

func (o *Square) SetA(a int) {
	o.a = a
}

func (o Square) GetArea() int {
	return o.a * o.a
}

func (o *Square) SetB(b int) {
	//
	// должно ли o.a быть равно b?
	// или оставаться неопределенным?
	//
}