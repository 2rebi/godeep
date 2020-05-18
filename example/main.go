package main

import (
	"fmt"
	"github.com/rebirthlee/godeep"
	"github.com/rebirthlee/godeep/example/unexport"
)

type Sample2 struct {
	Name string
	Action func() `from:"action"`
}

func main() {
	dst := Sample2{}
	godeep.Copy(&dst, unexport.New())
	fmt.Println(dst)
	dst.Action()
}
