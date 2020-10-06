package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testChar = uint16(0x41) // "A"

func TestLC3RAM_Load(t *testing.T) {
	t.Skip()
}

func TestLC3RAM_Read(t *testing.T) {

}

func TestLC3RAM_Write(t *testing.T) {
	m := &LC3RAM{
		CheckKey: KeyPressedMock(true),
		GetChar:  GetTestChar,
	}

	m.Write(0x100, uint16(0xFF))
	assert.Equal(t, uint16(0xFF), m.Read(0x100))

	address := m.Read(MR_KBSR)
	assert.Equal(t, "A", string(m.Storage[MR_KBDR]))
	assert.Equal(t, uint16(0b1000_0000_0000_0000), address)

	m.CheckKey = KeyPressedMock(false)
	address = m.Read(MR_KBSR)
	assert.Equal(t, uint16(0), address)
}

func KeyPressedMock(status bool) func() bool {
	return func() bool { return status }
}

func GetTestChar() uint16 {
	return testChar
}
