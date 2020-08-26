package instance

import "github.com/septemhill/test/recycle"

type instImpl struct {
	unit int
	grp  recycle.Group
}

func (i instImpl) Group() recycle.Group {
	return i.grp
}

func (i instImpl) Provided() int {
	return i.unit
}

//func NewInstance(g recycle.Group) recycle.Instance {
//	return &instImpl{
//		grp: g,
//	}
//}
