package unexport

import "fmt"

type sampple struct {
	age int
	Name string
	action func()
}


func New() *sampple {
	return &sampple{
		Name:   "Sample Name",
		action: func() {
			fmt.Println("hello world")
		},
	}
}