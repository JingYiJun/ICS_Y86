package models

import . "ICS_Y86_Backend/utils"

type Instruction struct {
	BitCode []byte
	LineNum uint64
	AsmCode string
	Comment string
}

func (i Instruction) ICode() uint64 {
	return Hex2uint64(i.BitCode[:1])
}

func (i Instruction) IFun() uint64 {
	return Hex2uint64(i.BitCode[1:2])
}

func (i Instruction) RegA() uint64 {
	return Hex2uint64(i.BitCode[2:3])
}

func (i Instruction) RegB() uint64 {
	return Hex2uint64(i.BitCode[3:4])
}

func (i Instruction) ValC() uint64 {
	var ans uint64 = 0
	for j := 0; j < 8; j++ {
		res := Hex2uint64(i.BitCode[2*j+4 : 2*j+6])
		ans = ans + res<<(j*8)
	}
	return ans
}

func (i Instruction) ValJmp() uint64 {
	var ans uint64 = 0
	for j := 0; j < 8; j++ {
		res := Hex2uint64(i.BitCode[2*j+2 : 2*j+4])
		ans = ans + res<<(j*8)
	}
	return ans
}

func (i Instruction) ParseAll() (iCode uint64, iFun uint64, RegA uint64, RegB uint64, ValC uint64, ValJmp uint64) {
	return i.ICode(), i.IFun(), i.RegA(), i.RegB(), i.ValC(), i.ValJmp()
}
