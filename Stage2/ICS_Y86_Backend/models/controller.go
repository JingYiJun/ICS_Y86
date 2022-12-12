package models

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Controller struct {
	Device  *Device
	Program map[uint64]Instruction
}

func NewController(RawCode *bufio.Reader) *Controller {
	c := &Controller{
		Device:  &Device{Stat: AOK},
		Program: map[uint64]Instruction{},
	}
	err := c.Parse(RawCode)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Controller) Parse(RawCode *bufio.Reader) error {
	scanner := bufio.NewScanner(RawCode)
	for scanner.Scan() {
		line := scanner.Bytes()

		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("|")) {
			continue
		}
		// cut line number
		addrBytes, remain, ok := bytes.Cut(line, []byte(":"))
		if !ok {
			return errors.New("finding line number failed")
		}
		addrBytes = bytes.TrimSpace(addrBytes)
		// parsing line number
		addr, err := strconv.ParseInt(string(addrBytes), 0, 64)
		if err != nil {
			return errors.New("parsing addr failed" + err.Error())
		}

		// cut instruction
		instruction, remain, ok := bytes.Cut(remain, []byte("|"))
		if !ok {
			return errors.New("finding instruction failed")
		}
		instruction = bytes.ToUpper(bytes.TrimSpace(instruction))
		if len(instruction) == 0 {
			continue
		}
		for _, c := range instruction {
			if !('0' <= c && c <= '9' || 'A' <= c && c <= 'F') {
				return errors.New("instruction text error")
			}
		}
		c.Device.Mem.WriteByte(uint64(addr), instruction)

		// cut asm code
		asmCode, comment, _ := bytes.Cut(remain, []byte("#"))
		asmCode = bytes.TrimSpace(asmCode)
		comment = bytes.TrimSpace(comment)

		// insert into Program
		c.Program[uint64(addr)] = Instruction{
			BitCode: bytes.Join([][]byte{instruction, bytes.Repeat([]byte{'0'}, 20-len(instruction))}, []byte{}),
			LineNum: uint64(addr),
			AsmCode: string(asmCode),
			Comment: string(comment),
		}
	}
	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return nil
}
