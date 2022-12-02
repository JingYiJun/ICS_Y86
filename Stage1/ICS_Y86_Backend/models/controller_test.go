package models

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestController_Parse(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(
		`
0x00a: 30f23f00000000000000 | irmovq $63, %rdx # src and dst have 63 elements
0x014: 30f69802000000000000 | irmovq dest, %rsi # dst array
0x01e: 30f79000000000000000 | irmovq src, %rdi # src array
`))
	c := NewController(reader)
	for _, v := range c.Program {
		fmt.Println(v)
	}
}
