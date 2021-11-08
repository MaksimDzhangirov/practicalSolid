package badCode4

type ConvexQuadrilateral interface {
	GetArea() int
}

type EquilateralRectangle interface {
	ConvexQuadrilateral
	SetA(a int)
}

type Oblong struct {
	EquilateralRectangle
	a int
	b int
}

func (o *Oblong) SetA(a int) {
	o.a = a
}

func (o *Oblong) SetB(b int) {
	// где определён этот метод?
	o.b = b
}

func (o Oblong) GetArea() int {
	return o.a * o.b
}

type Square struct {
	EquilateralRectangle
	a int
}

func (o *Square) SetA(a int) {
	o.a = a
}

func (o Square) GetArea() int {
	return o.a * o.a
}