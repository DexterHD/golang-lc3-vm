package main

import (
	"fmt"
	"os"
)

const PC_START uint16 = 0x3000

func main() {
	// {Load Arguments, 12}
	// {Setup, 12}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Missing argument!")
		return
	}

	vm := New(PC_START)
	vm.Load(args[0])
	vm.Run()
}
