package main

import (
	"fmt"

	"github.com/septemhill/tet/blah"
)

func main() {
	fmt.Println("Hi, Septem")

	g1 := blah.NewGroup("Dog", 100)
	g2 := blah.NewGroup("Cat", 100)

	g2.Digest(g1.Digested())
	g1.Digest(g2.Digested())
}
