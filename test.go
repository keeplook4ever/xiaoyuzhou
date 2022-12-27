package main

import (
	"fmt"
)

type Sss struct {
	Id int
}

func Gets() (Sss, error) {
	s := Sss{Id: 1}
	return s, nil
}


func main(){
	s, err := Gets()
	fmt.Println(s)
	fmt.Println(err)


}
