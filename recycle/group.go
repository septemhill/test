package recycle

import "time"

type Group interface {
	Count() int
	Digested() <-chan Instance
	Digest(<-chan Instance)
}

func genderStrategy() {
}

func noGenderStrategy() {
}

type grpImpl struct {
	grpName    string
	count      int
	inst       chan Instance
	digestList []<-chan Instance
	reproTime  time.Duration
	strategy   Reproduction
}

func (g grpImpl) Count() int {
	return g.count
}

func (g grpImpl) Digested() <-chan Instance {
	return g.inst
}

func (g grpImpl) Digest(ch <-chan Instance) {
	g.digestList = append(g.digestList, ch)
}

func (g grpImpl) reproduction() {

}

func NewGroup(name string, initCnt int, initRepro time.Duration, strategy Reproduction) Group {
	g := grpImpl{
		grpName:    name,
		inst:       make(chan Instance, 50),
		digestList: make([]<-chan Instance, 1),
		count:      initCnt,
		reproTime:  initRepro,
		strategy:   strategy,
	}

	return g
}
