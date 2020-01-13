package vm

import (
	"encoding/binary"
	"io/ioutil"
	"log"
)

// Memory
const MaxMemorySize uint16 = 65535

const (
	MR_KBSR uint16 = 0xfe00 // keyboard status
	MR_KBDR uint16 = 0xfe02 // keyboard data
)

type LC3RAM [MaxMemorySize]uint16

func (m *LC3RAM) Write(address, val uint16) {
	m[address] = val
}

func (m *LC3RAM) Read(address uint16) uint16 {
	if address == MR_KBSR {
		if checkKey() {
			m[MR_KBSR] = 1 << 15
			// read a single ASCII char
			m[MR_KBDR] = GetChar()
		} else {
			m[MR_KBSR] = 0
		}
	}
	return m[address]
}

func (m *LC3RAM) Load(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Can't read file, reason: %v", err)
	}

	origin := binary.BigEndian.Uint16(b[:2])

	for i := 2; i < len(b); i += 2 {
		m[origin] = binary.BigEndian.Uint16(b[i : i+2])
		origin++
	}
}
