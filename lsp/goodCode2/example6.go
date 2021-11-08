package goodCode2

type ConvexQuadrilateral interface {
	GetArea() int
}

type EquilateralQuadrilateral interface {
	ConvexQuadrilateral
	SetA(a int)
}

type NonEquilateralQuadrilateral interface {
	ConvexQuadrilateral
	SetA(a int)
	SetB(b int)
}

type NonEquiangularQuadrilateral interface {
	ConvexQuadrilateral
	SetAngle(angle float64)
}

type Oblong struct {
	NonEquilateralQuadrilateral
	a int
	b int
}

type Square struct {
	EquilateralQuadrilateral
	a int
}

type Parallelogram struct {
	NonEquilateralQuadrilateral
	NonEquiangularQuadrilateral
	a     int
	b     int
	angle float64
}

type Rhombus struct {
	EquilateralQuadrilateral
	NonEquiangularQuadrilateral
	a     int
	angle float64
}
