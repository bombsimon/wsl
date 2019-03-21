package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	flag.Usage = usage
	flag.Parse()

	if *html && *plumbing {
		log.Fatal("you can have either plumbing or HTML output")
	}
	if flag.NArg() > 0 {
		paths = flag.Args()
	}

	if *verbose {
		log.Println("Building suffix tree")
	}
	schan := job.Parse(filesFeed())
	t, data, done := job.BuildTree(schan)
	<-done

	fmt.Println(data)

	// finish stream
	t.Update(&syntax.Node{Type: -1})

	if *verbose {

		log.Println("Searching for clones")

	}

	somevar := 1
	somevar2 := 2
	if somevar {
		panic(somevar2)
	}

	_, err := thisFunc()
	if err != nil {
		ShouldBeOK()
	}

	if true {
		if false {
			if true {
				fmt.Println("We have to go deeper!")

			}
		}
	}
}

func thisFunc() (bool, error) {
	return false, nil
}

func ShouldBeOK() {
	fmt.Println("Don't be sad")
}
