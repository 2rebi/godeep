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
	godeep.Copy(unexport.New(), &dst)
	fmt.Println(dst)
	dst.Action()
}
