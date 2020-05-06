package main

import (
	"fmt"
	into "github.com/rebirthlee/golang-into"
)

type Sample1 struct {
	age int
	name string
}

type Sample2 struct {
	Number int `into:"age"`
	NickName string `into:"name"`
}

func main() {
	dst := Sample2{}
	fmt.Println(into.Into(Sample1{
		age: 10,
		name: "Sample",
	}, &dst))

	fmt.Println(dst)
}
