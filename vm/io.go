package vm

import (
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

// CheckKeyPushed checks if a key was pressed
func CheckKeyPressed() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Size() > 0
}

// GetCharFromStdin get one char from standard input.
func GetCharFromStdin() uint16 {
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
