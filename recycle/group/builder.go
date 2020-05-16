package group

import (
	"time"

	"github.com/septemhill/test/recycle"
)

type grpImplOpt func(g *grpImpl)

func GroupIsSexual(b bool) grpImplOpt {
	return func(g *grpImpl) {
		g.sexual = b
	}
}

func GroupStrategy(strategy recycle.Reproduction) grpImplOpt {
	return func(g *grpImpl) {
		g.strategy = strategy
	}
}

func GroupReproductionTime(t time.Duration) grpImplOpt {
	return func(g *grpImpl) {
		g.reproTicker = time.NewTicker(t)
	}
}

//func Group
