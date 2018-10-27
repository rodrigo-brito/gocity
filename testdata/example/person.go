package example

const BaseSalary = 10000

type Person struct {
	Name  string
	Years int
}

type Employee struct {
	Person
	Salary int
}

func NewEmployee(name string) Employee {
	person := Person{
		Name: name,
	}

	return Employee{
		Person: person,
		Salary: BaseSalary,
	}
}

func (e *Employee) SetYears(years int) {
	e.Years = years
}
