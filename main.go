package main

import (
	"fmt"
	"os"

	"github.com/idexter/golang-lc3-vm/vm"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Missing argument!")
		return
	}

	lc3 := vm.NewCPU(&vm.LC3RAM{
		CheckKey: vm.CheckKeyPressed,
		GetChar:  vm.GetCharFromStdin,
	}, os.Stdout)

	lc3.RAM.Load(args[0])
	lc3.Run()
}
