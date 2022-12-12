package models

import (
	. "ICS_Y86_Backend/utils"
	"bytes"
	"fmt"
)

type Device struct {
	PC   uint64  `json:"PC"`
	Reg  Regs    `json:"REG"`
	CC   CCs     `json:"CC"`
	Stat uint64  `json:"STAT"`
	Mem  *Memory `json:"MEM"`
}

type Regs []uint64
type CCs []uint64
type Memory [AddrMax]uint64

var RegName = [...]string{"rax", "rcx", "rdx", "rbx", "rsp", "rbp", "rsi", "rdi", "r8", "r9", "r10", "r11", "r12", "r13", "r14"}
var CCName = [...]string{"ZF", "SF", "OF"}
var MaxUsedMemory uint64 = 0

func (r Regs) MarshalJSON() ([]byte, error) {
	builder := bytes.Buffer{}
	builder.WriteByte('{')
	for i, v := range r {
		builder.WriteString(fmt.Sprintf("\"%v\":%v", RegName[i], int64(v)))
		if i != len(r)-1 {
			builder.WriteByte(',')
		}
	}
	builder.WriteByte('}')
	return builder.Bytes(), nil
}

func (c CCs) MarshalJSON() ([]byte, error) {
	builder := bytes.Buffer{}
	builder.WriteByte('{')
	for i, v := range c {
		builder.WriteString(fmt.Sprintf("\"%v\":%v", CCName[i], v))
		if i != len(c)-1 {
			builder.WriteByte(',')
		}
	}
	builder.WriteByte('}')
	return builder.Bytes(), nil
}

func (m *Memory) MarshalJSON() ([]byte, error) {
	builder := new(bytes.Buffer)
	builder.WriteByte('{')
	type MemIter struct{ k, v uint64 }
	tmp := new([]MemIter)
	for i := uint64(0); i < MaxUsedMemory+8; i += 8 {
		val, _ := m.Read(i, 8)
		if val != 0 {
			*tmp = append(*tmp, MemIter{k: i, v: val})
		}
	}
	for i, v := range *tmp {
		builder.WriteString(fmt.Sprintf("\"%v\":%v", v.k, int64(v.v)))
		if i != len(*tmp)-1 {
			builder.WriteByte(',')
		}
	}
	builder.WriteByte('}')
	return builder.Bytes(), nil
}

func (m *Memory) Write(addr, val, len uint64) uint64 {
	if addr >= AddrMax || addr+len >= AddrMax {
		return ADR
	}
	if addr > MaxUsedMemory {
		MaxUsedMemory = addr
	}
	for i := addr; i < addr+len; i++ {
		m[i] = val & 0xff
		val >>= 8
	}
	return AOK
}

func (m *Memory) WriteByte(addr uint64, val []byte) uint64 {
	l := uint64(len(val)) >> 1
	if addr+l >= AddrMax {
		return ADR
	}
	if addr+l > MaxUsedMemory {
		MaxUsedMemory = addr + l
	}
	for i := uint64(0); i < l; i++ {
		m[addr+i] = Hex2uint64(val[2*i : 2*i+2])
	}
	return AOK
}

func (m *Memory) Read(addr, len uint64) (uint64, uint64) {
	if addr >= AddrMax || addr+len >= AddrMax {
		return 0, ADR
	}
	var val uint64 = 0
	for i := uint64(0); i < len; i++ {
		val += m[addr+i] << (i * 8)
	}
	return val, AOK
}

func (d *Device) Push(val uint64) uint64 {
	d.Reg[4] -= 8
	return d.Mem.Write(d.Reg[4], val, 8)
}

func (d *Device) Pop() (uint64, uint64) {
	val, stat := d.Mem.Read(d.Reg[4], 8)
	d.Reg[4] += 8
	return val, stat
}

func (c CCs) CheckCondition(cond uint64) (bool, uint64) {
	switch cond {
	case 0:
		return true, AOK
	case 1: // jle (SF ^ OF) | ZF
		if (c[SF]^c[OF])|c[ZF] == 1 {
			return true, AOK
		} else {
			return false, AOK
		}
	case 2: // jl SF ^ OF
		if (c[SF] ^ c[OF]) == 1 {
			return true, AOK
		} else {
			return false, AOK
		}
	case 3: // je ZF
		if c[ZF] == 1 {
			return true, AOK
		} else {
			return false, AOK
		}
	case 4: // jne ~ZF
		if c[ZF] == 0 {
			return true, AOK
		} else {
			return false, AOK
		}
	case 5: // jge ~(SF ^ OF)
		if (c[SF] ^ c[OF]) == 0 {
			return true, AOK
		} else {
			return false, AOK
		}
	case 6: // jg ~((SF ^ OF) | ZF)
		if (c[SF]^c[OF])|c[ZF] == 0 {
			return true, AOK
		} else {
			return false, AOK
		}
	}
	return false, INS
}

func (d *Device) OP(iFun, valA, valB uint64) (uint64, uint64) {
	var ans int64 = 0
	a := int64(valA)
	b := int64(valB)

	switch iFun {

	case 0: // addq rA, rB
		ans = b + a
		if ((a < 0) == (b < 0)) && ((ans < 0) != (a < 0)) {
			d.CC[OF] = 1
		} else {
			d.CC[OF] = 0
		}

	case 1: // subq rA, rB
		ans = b - a
		if ((a > 0) == (b < 0)) && ((ans < 0) != (a > 0)) {
			d.CC[OF] = 1
		} else {
			d.CC[OF] = 0
		}

	case 2: // andq rA, rB
		ans = b & a

	case 3: // xorq rA, rB
		ans = b ^ a

	default:
		return uint64(b), INS
	}

	if ans == 0 {
		d.CC[ZF] = 1
	} else {
		d.CC[ZF] = 0
	}

	if uint64(ans)&(1<<63) != 0 {
		d.CC[SF] = 1
	} else {
		d.CC[SF] = 0
	}
	return uint64(ans), AOK
}
