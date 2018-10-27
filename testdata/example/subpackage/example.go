package subpackage

type AnimalType int

const (
	Domestic AnimalType = iota
	Wild
)

const AnimalFlag = iota

type Animal struct {
	Name   string
	Specie string
	wild   bool
}

func NewAnimal(name string, flag AnimalType) Animal {
	animal := Animal{
		Name: name,
	}

	switch flag {
	case Domestic:
		animal.wild = true
	default:
		animal.wild = false
	}

	return animal
}
