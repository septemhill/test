package recycle

type Reproduction interface {
	Reproduce()
}

type GenderStrategy struct{}

func (s GenderStrategy) Reproduce() {}

type NoGenderStrategy struct{}

func (s NoGenderStrategy) Reproduce() {}
