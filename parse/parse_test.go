package parse_test

import (
	"fmt"
	"github.com/x0y14/beat/parse"
	"github.com/x0y14/beat/tokenize"
	"testing"
)

func TestParseWork(t *testing.T) {
	code := `
	var gA int = 1
	func sayHelloS(name string) string {
		var hello string = "hello, "
		return hello + name + "!"
	}
	func sayHello(name string) {
		fmt.printf(sayHelloS(name))
	}
	func sub(x int, y int) (int, bool) {
		var isMinus bool
		z := x - y
		if z < 0 {
			isMinus = true
		} else {
			isMinus = false
		}
		return z, isMinus
	}
	`
	head, err := tokenize.Tokenize(code)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := parse.Parse(head)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nodes)
}
