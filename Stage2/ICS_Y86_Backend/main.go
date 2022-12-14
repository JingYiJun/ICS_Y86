//go:build wasm

package main

import (
	"ICS_Y86_Backend/models"
	"bufio"
	"io"
	"strings"
	"syscall/js"
)

type Map = map[string]any

func run(reader io.Reader) Map {
	r := models.Controller{
		Device: &models.Device{
			PC:   0,
			Reg:  make([]uint64, 15, 15),
			CC:   []uint64{1, 0, 0},
			Stat: models.AOK,
			Mem:  &models.Memory{},
		},
		Program: map[uint64]models.Instruction{},
	}
	err := r.Parse(bufio.NewReader(reader))
	if err != nil {
		return Map{"result": "", "error": err.Error()}
	}
	result, err := r.Run()
	if err != nil {
		return Map{"result": "", "error": err.Error()}
	}
	return Map{"result": result, "error": ""}
}

func goRun(_ js.Value, args []js.Value) interface{} {
	reader := strings.NewReader(args[0].String())
	return js.ValueOf(run(reader))
}

func main() {
	done := make(chan int, 0)
	js.Global().Set("goRun", js.FuncOf(goRun))
	<-done
}
