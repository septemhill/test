package recycle

type Reproduction interface {
	Reproduce() int
}

type GenderStrategy struct{}

func (s GenderStrategy) Reproduce() int {
	return 1000
}

type NoGenderStrategy struct{}

func (s NoGenderStrategy) Reproduce() int {
	return 5000
}
