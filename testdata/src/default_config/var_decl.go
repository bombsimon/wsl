package testpkg

import (
	"net/http"
	"time"
)


func fn() {
    var a = 1 // want "declarations should never be cuddled"
    var baba  = 2 // what about
    var xx  = 2 // Some comments?
    var yy  = 2

    var c = 1 // want "declarations should never be cuddled"
	var (
		d             = 1 // like this
		extendedAlign = 33
	)
	var (
		x = 1
		y = 2
	)
	var afsa = 1

	if true {
		//
	}
	var aa = "a" // want "declarations should never be cuddled"
	var bb  = "b"

	var s = http.Server{ // want "declarations should never be cuddled"
		Addr:        ":8080", // important comment
		ReadTimeout: time.Second,
	}
	var s2 = http.Server{
		Addr:        ":8081",
		ReadTimeout: 2 * time.Second,
	}

	var a1 = 1 // want "declarations should never be cuddled"
	var s1 = http.Server{
		Addr:        ":8080", // important comment
		ReadTimeout: time.Second,
	}

	var s2 = http.Server{ // want "declarations should never be cuddled"
		Addr:        ":8080", // important comment
		ReadTimeout: time.Second,
	}
	var a2 = 1 // also comment

	var s = http.Server{
		Addr:        ":8080", // important comment
		ReadTimeout: time.Second,
	}
	type T int // want "declarations should never be cuddled"

	var a int // want "declarations should never be cuddled"
	var b int
	type T int // want "declarations should never be cuddled"

	for _, n := range []int{1, 2, 3} {
		var s *http.Server // want "declarations should never be cuddled"
		var err error
		switch n { // want "only one cuddle assignment allowed before switch statement"
		case 1:
			fmt.Println(err)
		case 2:
			fmt.Println(s)
		}
	}

	// Some comment
	const a = 1
	var b = 2 // want "declarations should never be cuddled"
}
