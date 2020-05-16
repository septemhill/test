package recycle

type Instance interface {
	Type() string
	Provided() int
}

type instImpl struct {
	typ  string
	unit int
}

func (i instImpl) Type() string {
	return i.typ
}

func (i instImpl) Provided() int {
	return i.unit
}
