package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVirtualMachine_updateFlags(t *testing.T) {
	vm := NewCpu()
	vm.registers[R_R0] = 0b0000000000000001 // int16(1)
	vm.registers[R_R1] = 0b1111111111111111 // int16(-1)
	vm.registers[R_R2] = 0b0000000000000000 // int16(0)

	vm.updateFlags(R_R0)
	assert.Equal(t, FL_POS, vm.registers[R_COND])
	vm.updateFlags(R_R1)
	assert.Equal(t, FL_NEG, vm.registers[R_COND])
	vm.updateFlags(R_R2)
	assert.Equal(t, FL_ZRO, vm.registers[R_COND])
}

func TestVirtualMachine_add(t *testing.T) {
	vm := NewCpu()

	vm.registers[R_R1] = 0x1
	vm.registers[R_R2] = 0x2
	vm.currentInstruction = 0b0001_000_001_0_00_010 // ADD R_R0, R_R1, R_R2
	vm.add()

	assert.Equal(t, uint16(0x3), vm.registers[R_R0])
	assert.Equal(t, uint16(0x1), vm.registers[R_R1])
	assert.Equal(t, uint16(0x2), vm.registers[R_R2])
	vm.Reset()

	vm.registers[R_R1] = 0x4
	vm.currentInstruction = 0b0001_000_001_1_00110 // ADD R_R0, R_R1, 6
	vm.add()

	assert.Equal(t, uint16(0xA), vm.registers[R_R0])
	vm.Reset()
}

func TestVirtualMachine_and(t *testing.T) {
	vm := NewCpu()

	vm.registers[R_R1] = 0x9
	vm.registers[R_R2] = 0x7
	vm.currentInstruction = 0b0001_000_001_0_00_010 // AND R_R0, R_R1, R_R2
	vm.and()

	assert.Equal(t, uint16(0x1), vm.registers[R_R0])
	assert.Equal(t, uint16(0x9), vm.registers[R_R1])
	assert.Equal(t, uint16(0x7), vm.registers[R_R2])
	vm.Reset()

	vm.registers[R_R1] = 0x9
	vm.currentInstruction = 0b0001_000_001_1_00111 // AND R_R0, R_R1, 7
	vm.and()

	assert.Equal(t, uint16(0x1), vm.registers[R_R0])
	vm.Reset()
}

func TestVirtualMachine_not(t *testing.T) {
	vm := NewCpu()

	vm.registers[R_R0] = 0xffff
	vm.registers[R_R1] = 0x0000
	vm.currentInstruction = 0b1001_000_001_000000 // NOT R_R0, R_R1
	vm.not()

	assert.Equal(t, uint16(0xffff), vm.registers[R_R0])
	assert.Equal(t, uint16(0x0000), vm.registers[R_R1])
	vm.Reset()
}

func TestVirtualMachine_branch(t *testing.T) {
	vm := NewCpu()

	vm.RAM.Storage[0x0003] = 0x1111
	vm.registers[R_COND] = FL_ZRO
	vm.currentInstruction = 0b0000_010_000000011

	vm.branch()

	assert.Equal(t, uint16(0x1111), vm.RAM.Storage[vm.registers[R_PC]])
	vm.Reset()

	vm.RAM.Storage[0x0003] = 0x2222
	vm.registers[R_COND] = FL_NEG
	vm.currentInstruction = 0b0000_100_000000011

	vm.branch()

	assert.Equal(t, uint16(0x2222), vm.RAM.Storage[vm.registers[R_PC]])
	vm.Reset()

	vm.RAM.Storage[0x0003] = 0x3333
	vm.registers[R_COND] = FL_POS
	vm.currentInstruction = 0b0000_001_000000011

	vm.branch()

	assert.Equal(t, uint16(0x3333), vm.RAM.Storage[vm.registers[R_PC]])
	vm.Reset()
}

func TestLC3CPU_jump(t *testing.T) {
	vm := NewCpu()

	vm.RAM.Storage[0x0001] = 0x1111
	vm.registers[R_R2] = 0x3100
	vm.currentInstruction = 0b1100_000_010_000000

	vm.jump()

	assert.Equal(t, uint16(0x3100), vm.registers[R_PC])
	vm.Reset()
}

func TestLC3CPU_jumpRegister(t *testing.T) {
	vm := NewCpu()

	vm.registers[R_PC] = PC_START
	vm.currentInstruction = 0b0100_1_00000001000

	vm.jumpRegister()

	assert.Equal(t, uint16(0x3008), vm.registers[R_PC])
	vm.Reset()

	vm.registers[R_PC] = PC_START
	vm.registers[R_R2] = 0x3100
	vm.currentInstruction = 0b0100_0_00_010_000000

	vm.jumpRegister()

	assert.Equal(t, uint16(0x3100), vm.registers[R_PC])
	vm.Reset()
}

func TestLC3CPU_load(t *testing.T) {
	vm := NewCpu()

	vm.RAM.Write(0x3004, 2)

	vm.registers[R_PC] = PC_START
	vm.currentInstruction = 0b0010_010_000000100
	vm.load()

	assert.Equal(t, uint16(2), vm.registers[R_R2])
	vm.Reset()
}

func TestLC3CPU_ldi(t *testing.T) {
	vm := NewCpu()

	vm.RAM.Write(0x3004, 0x3008)
	vm.RAM.Write(0x3008, 5)

	vm.registers[R_PC] = PC_START
	vm.currentInstruction = 0b1010_010_000000100
	vm.ldi()

	assert.Equal(t, uint16(5), vm.registers[R_R2])
	vm.Reset()
}

func TestLC3CPU_loadRegister(t *testing.T) {
	vm := NewCpu()

	vm.RAM.Write(0x3014, 5)

	vm.registers[R_PC] = PC_START
	vm.registers[R_R4] = 0x3010
	vm.currentInstruction = 0b0110_010_100_000100
	vm.loadRegister()

	assert.Equal(t, uint16(5), vm.registers[R_R2])
	vm.Reset()
}

func TestLC3CPU_loadEffectiveAddress(t *testing.T) {
	vm := NewCpu()

	vm.RAM.Write(0x3014, 5)

	vm.registers[R_PC] = PC_START
	vm.currentInstruction = 0b1110_010_000000100
	vm.loadEffectiveAddress()

	assert.Equal(t, uint16(0x3004), vm.registers[R_R2])
	vm.Reset()
}

func TestLC3CPU_store(t *testing.T) {
	vm := NewCpu()

	vm.registers[R_PC] = PC_START
	vm.registers[R_R2] = 0xff
	vm.currentInstruction = 0b0011_010_000000100
	vm.store()

	assert.Equal(t, uint16(0xff), vm.RAM.Read(0x3004))
	vm.Reset()
}

func TestLC3CPU_storeIndirect(t *testing.T) {
	vm := NewCpu()

	vm.RAM.Write(0x3004, 0x3010)

	vm.registers[R_PC] = PC_START
	vm.registers[R_R2] = 0xff
	vm.currentInstruction = 0b1011_010_000000100
	vm.storeIndirect()

	assert.Equal(t, uint16(0xff), vm.RAM.Read(0x3010))
	vm.Reset()
}

func TestLC3CPU_storeRegister(t *testing.T) {
	vm := NewCpu()

	vm.registers[R_PC] = PC_START
	vm.registers[R_R1] = 0x3001
	vm.registers[R_R2] = 0xff
	vm.currentInstruction = 0b0111_010_001_000100
	vm.storeRegister()

	assert.Equal(t, uint16(0xff), vm.RAM.Read(0x3005))
	vm.Reset()
}
