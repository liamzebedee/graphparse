package testpkg

type Church struct {
	balance int
	believers []Person
	pope Person
}

func NewChurch(balance int) Church {
	return Church{}
}

func (church *Church) AddBeliever(p Person) {
	church.believers = append(church.believers, p)
}

func (church *Church) AnnounceDayOfOurLord() {
	church.balance = 10000
}

type Person struct {
	name string
}

func NewPerson(name string) Person {
	return Person{
		name: name,
	}
}