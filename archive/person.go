package archive

import "github.com/genealogix/glx/go-glx"

type Person struct {
	g *glx.Person
}

func newPerson(g *glx.Person) *Person {
	return new(Person{g: g})
}
