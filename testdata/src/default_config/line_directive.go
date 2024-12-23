//line line_directive.tmpl:1
package pkg

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
)

func _() {
	var x int
	_ = x
	x = 0 // x
}

func main() {
	a()
	b()
	wsl()
}

func a() {
	fmt.Println("foo")
}

func b() {
	fmt.Println("foo")
}

func c() {
	_ = analysis.Analyzer{}
}

// want +1 "block should not start with a whitespace"
func wsl() bool {

	return true
}

func notFormatted() {
}
