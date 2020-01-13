package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

// Memory
const UINT16_MAX uint16 = 65535

const (
	MR_KBSR uint16 = 0xfe00 // keyboard status
	MR_KBDR uint16 = 0xfe02 // keyboard data
)

// Registers
const (
	R_R0 uint16 = iota
	R_R1
	R_R2
	R_R3
	R_R4
	R_R5
	R_R6
	R_R7
	R_PC // program counter
	R_COND
	R_COUNT
)

// Opcodes
const (
	OP_BR   uint16 = iota // branch
	OP_ADD                // add
	OP_LD                 // load
	OP_ST                 // store
	OP_JSR                // jump register
	OP_AND                // bitwise and
	OP_LDR                // load register
	OP_STR                // store register
	OP_RTI                // unused
	OP_NOT                // bitwise not
	OP_LDI                // load indirect
	OP_STI                // store indirect
	OP_JMP                // jump
	OP_RES                // reserved (unused)
	OP_LEA                // load effective address
	OP_TRAP               // execute trap
)

// Condition Flags
const (
	FL_POS uint16 = 1 << 0 // Positive
	FL_ZRO uint16 = 1 << 1 // Zero
	FL_NEG uint16 = 1 << 2 // Negative
)

type VirtualMachine struct {
	registers     [R_COUNT]uint16
	memory        [UINT16_MAX]uint16
	instruction   uint16
	operation     uint16
	running       bool
	startPosition uint16
}

func New(startPosition uint16) *VirtualMachine {
	return &VirtualMachine{
		startPosition: startPosition,
	}
}

func (v *VirtualMachine) Load(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Can't read file, reason: %v", err)
	}

	origin := binary.BigEndian.Uint16(b[:2])

	for i := 2; i < len(b); i += 2 {
		v.memory[origin] = binary.BigEndian.Uint16(b[i : i+2])
		origin++
	}
}

func (v *VirtualMachine) Reset() {
	v.registers = [R_COUNT]uint16{}
	v.memory = [UINT16_MAX]uint16{}
	v.instruction = 0
	v.operation = 0
	v.running = false
}

func (v *VirtualMachine) Run() {
	// Set the PC to starting position
	// 0x3000 is the default
	v.registers[R_PC] = v.startPosition
	v.running = true
	for v.running {
		// Fetch
		v.instruction = v.ReadFromMemory(v.registers[R_PC])
		if v.registers[R_PC] < UINT16_MAX {
			v.registers[R_PC]++
		}
		v.operation = v.instruction >> 12

		switch v.operation {
		case OP_ADD:
			v.add()
			break
		case OP_AND:
			v.and()
			break
		case OP_NOT:
			v.not()
			break
		case OP_BR:
			v.branch()
			break
		case OP_JMP:
			v.jump()
			break
		case OP_JSR:
			v.jumpRegister()
			break
		case OP_LD:
			v.load()
			break
		case OP_LDI:
			v.ldi()
			break
		case OP_LDR:
			v.loadRegister()
			break
		case OP_LEA:
			v.loadEffectiveAddress()
			break
		case OP_ST:
			v.store()
			break
		case OP_STI:
			v.storeIndirect()
			break
		case OP_STR:
			v.storeRegister()
			break
		case OP_TRAP:
			switch v.instruction & 0xFF {
			case TRAP_GETC:
				v.trapGetc()
				break
			case TRAP_OUT:
				v.trapOut()
				break
			case TRAP_PUTS:
				v.trapPuts()
				break
			case TRAP_IN:
				v.trapIn()
				break
			case TRAP_PUTSP:
				v.trapPutsp()
				break
			case TRAP_HALT:
				v.trapHalt()
				break
			}
			break
		case OP_RES:
		case OP_RTI:
		default:
			log.Printf("BAD OPCODE: %016b", v.operation)
			v.running = false
		}
	}
}

func (v *VirtualMachine) updateFlags(r uint16) {
	if v.registers[r] == 0 {
		v.registers[R_COND] = FL_ZRO
	} else if v.registers[r]>>15 == uint16(1) { //* a 1 in the left-most bit indicates negative */
		v.registers[R_COND] = FL_NEG
	} else {
		v.registers[R_COND] = FL_POS
	}
}

func (v *VirtualMachine) WriteToMemory(address, val uint16) {
	v.memory[address] = val
}

func (v *VirtualMachine) ReadFromMemory(address uint16) uint16 {
	if address == MR_KBSR {
		if checkKey() {
			v.memory[MR_KBSR] = 1 << 15
			// read a single ASCII char
			v.memory[MR_KBDR] = GetChar()
		} else {
			v.memory[MR_KBSR] = 0
		}
	}
	return v.memory[address]
}

// -------------- Instruction Implementations --------------------

func (v *VirtualMachine) add() {
	// destination register (DR)
	r0 := (v.instruction >> 9) & 0x7
	// first operand (SR1)
	r1 := (v.instruction >> 6) & 0x7
	// whether we are in immediate mode
	immFlag := (v.instruction >> 5) & 0x1

	if immFlag == 0x1 {
		imm5 := SignExtend(v.instruction&0x1F, 5)
		v.registers[r0] = v.registers[r1] + imm5
	} else {
		r2 := v.instruction & 0x7
		v.registers[r0] = v.registers[r1] + v.registers[r2]
	}

	v.updateFlags(r0)
}

func (v *VirtualMachine) and() {
	r0 := (v.instruction >> 9) & 0x7
	r1 := (v.instruction >> 6) & 0x7
	immFlag := (v.instruction >> 5) & 0x1

	if immFlag == 0x1 {
		imm5 := SignExtend(v.instruction&0x1F, 5)
		v.registers[r0] = v.registers[r1] & imm5
	} else {
		r2 := v.instruction & 0x7
		v.registers[r0] = v.registers[r1] & v.registers[r2]
	}
	v.updateFlags(r0)
}

func (v *VirtualMachine) not() {
	r0 := (v.instruction >> 9) & 0x7
	r1 := (v.instruction >> 6) & 0x7

	v.registers[r0] = ^v.registers[r1]
	v.updateFlags(r0)
}

func (v *VirtualMachine) branch() {
	pcOffset := SignExtend((v.instruction)&0x1ff, 9)
	condFlag := (v.instruction >> 9) & 0x7
	if (condFlag & v.registers[R_COND]) != 0 { // true
		v.registers[R_PC] += pcOffset
	}
}

func (v *VirtualMachine) jump() {
	/* Also handles RET */
	r1 := (v.instruction >> 6) & 0x7
	v.registers[R_PC] = v.registers[r1]
}

func (v *VirtualMachine) jumpRegister() {
	r1 := (v.instruction >> 6) & 0x7
	longPcOffset := SignExtend(v.instruction&0x7ff, 11)
	longFlag := (v.instruction >> 11) & 1

	v.registers[R_R7] = v.registers[R_PC]
	if longFlag == 1 {
		v.registers[R_PC] += longPcOffset /* JSR */
	} else {
		v.registers[R_PC] = v.registers[r1] /* JSRR */
	}
}

func (v *VirtualMachine) load() {
	r0 := (v.instruction >> 9) & 0x7
	pcOffset := SignExtend(v.instruction&0x1ff, 9)
	v.registers[r0] = v.ReadFromMemory(v.registers[R_PC] + pcOffset)
	v.updateFlags(r0)
}

func (v *VirtualMachine) ldi() {
	/* destination register (DR) */
	r0 := (v.instruction >> 9) & 0x7
	/* PCoffset 9*/
	pcOffset := SignExtend(v.instruction&0x1ff, 9)
	/* add pcOffset to the current PC, look at that memory location to get the final address */
	v.registers[r0] = v.ReadFromMemory(v.ReadFromMemory(v.registers[R_PC] + pcOffset))
	v.updateFlags(r0)
}

func (v *VirtualMachine) loadRegister() {
	r0 := (v.instruction >> 9) & 0x7
	r1 := (v.instruction >> 6) & 0x7
	offset := SignExtend(v.instruction&0x3F, 6)
	v.registers[r0] = v.ReadFromMemory(v.registers[r1] + offset)
	v.updateFlags(r0)
}

func (v *VirtualMachine) loadEffectiveAddress() {
	r0 := (v.instruction >> 9) & 0x7
	pcOffset := SignExtend(v.instruction&0x1ff, 9)
	v.registers[r0] = v.registers[R_PC] + pcOffset
	v.updateFlags(r0)
}

func (v *VirtualMachine) store() {
	r0 := (v.instruction >> 9) & 0x7
	pcOffset := SignExtend(v.instruction&0x1ff, 9)
	v.WriteToMemory(v.registers[R_PC]+pcOffset, v.registers[r0])
}

func (v *VirtualMachine) storeIndirect() {
	r0 := (v.instruction >> 9) & 0x7
	pcOffset := SignExtend(v.instruction&0x1ff, 9)
	v.WriteToMemory(v.ReadFromMemory(v.registers[R_PC]+pcOffset), v.registers[r0])
}

func (v *VirtualMachine) storeRegister() {
	r0 := (v.instruction >> 9) & 0x7
	r1 := (v.instruction >> 6) & 0x7
	offset := SignExtend(v.instruction&0x3F, 6)
	v.WriteToMemory(v.registers[r1]+offset, v.registers[r0])
}

const (
	TRAP_GETC  = 0x20 // get character from keyboard, not echoed onto the terminal
	TRAP_OUT   = 0x21 // output a character
	TRAP_PUTS  = 0x22 // output a word string
	TRAP_IN    = 0x23 // get character from keyboard, echoed onto the terminal
	TRAP_PUTSP = 0x24 // output a byte string
	TRAP_HALT  = 0x25 // halt the program
)

func (v *VirtualMachine) trapGetc() {
	// read a single ASCII char
	v.registers[R_R0] = GetChar()
}

func (v *VirtualMachine) trapOut() {
	fmt.Printf("%c", v.registers[R_R0])
}

func (v *VirtualMachine) trapPuts() {
	for i := v.registers[R_R0]; v.memory[i] != 0x0000; i++ {
		fmt.Printf("%c", v.memory[i])
	}
	fmt.Print("\n")
}

func (v *VirtualMachine) trapIn() {
	fmt.Print("Input a character: ")
	c := GetChar()
	fmt.Printf("%c", c)
	v.registers[R_R0] = c
}

func (v *VirtualMachine) trapPutsp() {
	// one char per byte (two bytes per word)
	// here we need to swap back to
	// big endian format
	for i := v.registers[R_R0]; v.memory[i] > 0; i++ {
		ch1 := v.memory[i] & 0xFF
		fmt.Printf("%c", ch1)
		ch2 := v.memory[i] >> 8
		if ch2 > 0 {
			fmt.Printf("%c", ch2)
		}
	}
	fmt.Print("\n")
}

func (v *VirtualMachine) trapHalt() {
	fmt.Println("HALT")
	v.running = false
}
