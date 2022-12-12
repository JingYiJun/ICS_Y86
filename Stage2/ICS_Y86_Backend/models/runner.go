package models

import (
	"encoding/json"
	"strings"
)

const (
	AOK = iota + 1
	HLT
	ADR
	INS
)

const AddrMax = 64 * 1024 * 1024 // 64MB

func (c *Controller) Run() string {
	d := c.Device
	output := make([]string, 0, 100)
	for c.Device.Stat == AOK {
		// fetch
		instruction, ok := c.Program[d.PC]
		oldPC := d.PC
		if !ok {
			d.Stat = ADR
			continue
		}

		iCode, iFun, regA, regB, valC, valJmp := instruction.ParseAll()

		switch iCode {

		case 0: // halt
			switch iFun {
			case 0:
				d.Stat = HLT
			default:
				d.Stat = INS
			}

		case 1: // nop
			d.PC += 1
			switch iFun {
			case 0:
				goto OUTPUT
			default:
				d.Stat = INS
			}

		case 2: // rrmovq rA, rB
			d.PC += 2
			if regA == 15 || regB == 15 {
				d.Stat = INS
				goto OUTPUT
			}
			if ok, d.Stat = d.CC.CheckCondition(iFun); ok && d.Stat == AOK {
				d.Reg[regB] = d.Reg[regA]
			}

		case 3: // irmovq V, rB
			d.PC += 10
			switch iFun {
			case 0:
				if regA != 15 || regB == 15 {
					d.Stat = INS
				} else {
					d.Reg[regB] = valC
				}

			default:
				d.Stat = INS
			}

		case 4: // rmmovq rA, D(rB)
			d.PC += 10
			switch iFun {
			case 0:
				if regA == 15 || regB == 15 {
					d.Stat = INS
					goto OUTPUT
				}
				val := d.Reg[regA]
				addr := d.Reg[regB] + valC
				d.Stat = d.Mem.Write(addr, val, 8)

			default:
				d.Stat = INS
			}

		case 5: // mrmovq D(rB), rA
			d.PC += 10
			switch iFun {
			case 0:
				if regA == 15 || regB == 15 {
					d.Stat = INS
					goto OUTPUT
				}
				addr := d.Reg[regB] + valC
				d.Reg[regA], d.Stat = d.Mem.Read(addr, 8)

			default:
				d.Stat = INS
			}

		case 6: // OPq rA, rB
			d.PC += 2
			if regA == 15 || regB == 15 {
				d.Stat = INS
				goto OUTPUT
			}
			d.Reg[regB], d.Stat = d.OP(iFun, d.Reg[regA], d.Reg[regB])

		case 7: // jxx Dest
			d.PC += 9
			if valJmp > AddrMax {
				d.Stat = ADR
				goto OUTPUT
			}

			if ok, d.Stat = d.CC.CheckCondition(iFun); ok && d.Stat == AOK {
				d.PC = valJmp
			}

		case 8: // call Dest
			d.PC += 9
			switch iFun {
			case 0:
				if valJmp > AddrMax {
					d.Stat = ADR
					goto OUTPUT
				}
				d.Push(d.PC)
				d.PC = valJmp
			default:
				d.Stat = INS
			}

		case 9: // ret
			d.PC += 1
			if iFun != 0 {
				d.Stat = INS
				goto OUTPUT
			}

			if valJmp > AddrMax {
				d.Stat = ADR
				goto OUTPUT
			}
			d.PC, d.Stat = d.Pop()

		case 10: // pushq rA
			d.PC += 2
			if iFun != 0 || regA == 15 || regB != 15 {
				d.Stat = INS
				goto OUTPUT
			}
			d.Stat = d.Push(d.Reg[regA])

		case 11:
			d.PC += 2
			if iFun != 0 || regA == 15 || regB != 15 {
				d.Stat = INS
				goto OUTPUT
			}
			d.Reg[regA], d.Stat = d.Pop()

		case 12:
			d.PC += 10
			if regA != 15 || regB == 15 {
				d.Stat = INS
				goto OUTPUT
			}
			d.Reg[regB], d.Stat = d.OP(iFun, d.Reg[regB], valC)

		default:
			d.Stat = INS
		}

	OUTPUT:
		if d.Stat != AOK {
			d.PC = oldPC
		}
		ans, err := json.Marshal(d)
		if err != nil {
			panic(err)
		}
		output = append(output, string(ans))
	}

	outputBuilder := strings.Builder{}
	outputBuilder.WriteByte('[')
	for i, v := range output {
		outputBuilder.WriteString(v)
		if i != len(output)-1 {
			outputBuilder.WriteByte(',')
		}
	}
	outputBuilder.WriteByte(']')
	return outputBuilder.String()
}
