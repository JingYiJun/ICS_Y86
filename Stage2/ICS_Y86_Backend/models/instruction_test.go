package models

import (
	"fmt"
	"testing"
)

func TestInstructionParsing(t *testing.T) {
	instructions := []Instruction{
		{
			BitCode: []byte("30f23f3f3f3f3f3f3f3f"),
		},
		{
			BitCode: []byte("30f69802000000000000"),
		},
		{
			BitCode: []byte("30f79000000000000000"),
		},
	}

	for _, instruction := range instructions {
		fmt.Printf(
			`Code: %v,
ICode: %v,
IFun: %v,
RegA: %v,
RegB: %v,
ValV: %v,
ValJmp: %v,
`,
			string(instruction.BitCode),
			instruction.ICode(),
			instruction.IFun(),
			instruction.RegA(),
			instruction.RegB(),
			instruction.ValC(),
			instruction.ValJmp(),
		)
	}

}
