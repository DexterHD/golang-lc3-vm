package main

import (
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func SignExtend(x uint16, bitCount int) uint16 {
	if (x>>(bitCount-1))&1 == 1 {
		x |= 0xFFFF << bitCount
	}
	return x
}

// checkKey checks if a key was pressed
func checkKey() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Size() > 0
}

func GetChar() uint16 {
	// fd 0 is stdin
	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}
	defer func() {
		if err := terminal.Restore(0, state); err != nil {
			log.Println("warning, failed to restore terminal:", err)
		}
	}()

	b := make([]byte, 1)
	os.Stdin.Read(b)
	return uint16(b[0])
}
