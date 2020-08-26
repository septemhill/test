package recycle

type Cycle interface {
	Groups() []Group
}

type Group interface {
	Name() string
	Sexual() bool
	Count() int
	Digested() <-chan Instance
	Digest(<-chan Instance)
	Cycle() []Cycle
}

type Instance interface {
	Provided() int
	Group() Group
}
