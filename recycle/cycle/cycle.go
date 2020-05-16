package cycle

import "github.com/septemhill/test/recycle"

type cycleImpl struct {
	grps []recycle.Group
}

func (c cycleImpl) Groups() []recycle.Group {
	return c.grps
}

func NewCycle(grps ...recycle.Group) recycle.Cycle {
	if len(grps) <= 1 {
		return nil
	}

	for i := 0; i < len(grps); i++ {
		instCh := grps[i].Digested()
		grps[i+1].Digest(instCh)
	}

	grps[0].Digest(grps[len(grps)-1].Digested())

	c := cycleImpl{
		grps: grps,
	}

	return c
}
