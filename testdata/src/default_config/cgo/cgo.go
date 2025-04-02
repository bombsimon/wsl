package testpkg

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func Advanced() {
	var foo = 1
	var bar = 2
	var biz int

	x := []string{}
	x = append(x, "literal")
	notUsed := "just assigning, don't mind me"
	x = append(x, "not z..")
	useMe := "right away"
	alsoNotUsed := ":("
	x = append(x, "just noise")
	x = append(x, useMe)

	_ = foo
	_ = bar
	_ = biz
	_ = notUsed
	_ = alsoNotUsed
}
