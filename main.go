package main

import (
	"fmt"
	"laundry-room/pkg/laundry"
	"laundry-room/pkg/loader"
	"os"
)

func main(){
	a := laundry.NewLaundryApp()
	l := loader.NewLoader(a)
	if err := l.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}