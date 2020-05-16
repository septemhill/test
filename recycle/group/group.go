package group

import (
	"time"

	"github.com/septemhill/test/recycle"
)

type grpImpl struct {
	grpName     string
	maleCount   int
	femaleCount int
	sexual      bool
	count       int
	inst        chan recycle.Instance
	digestList  []<-chan recycle.Instance
	reproTicker *time.Ticker
	strategy    recycle.Reproduction
	cycle       []recycle.Cycle
}

func (g *grpImpl) Sexual() bool {
	return g.sexual
}

func (g *grpImpl) Count() int {
	return g.count
}

func (g *grpImpl) Digested() <-chan recycle.Instance {
	return g.inst
}

func (g *grpImpl) Digest(ch <-chan recycle.Instance) {
	g.digestList = append(g.digestList, ch)
}

func (g *grpImpl) reproduction() {
	if g.strategy == nil {
		return
	}

	g.count += g.strategy.Reproduce()
}

func (g *grpImpl) Name() string {
	return g.grpName
}

func (g *grpImpl) Cycle() []recycle.Cycle {
	return g.cycle
}

func (g *grpImpl) run() {
	for {
		select {
		case <-g.reproTicker.C:
			g.reproduction()
		}
	}
}

func NewGroup(name string, initCnt int, opts ...grpImplOpt) recycle.Group {
	g := &grpImpl{
		grpName:     name,
		inst:        make(chan recycle.Instance, 50),
		digestList:  make([]<-chan recycle.Instance, 1),
		count:       initCnt,
		reproTicker: time.NewTicker(time.Second * 2),
		strategy:    recycle.GenderStrategy{},
		cycle:       make([]recycle.Cycle, 0),
	}

	for _, opt := range opts {
		opt(g)
	}

	go g.run()

	return g
}
