package main

import (
	"ICS_Y86_Backend/models"
	"bufio"
	"os"
)

func main() {
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
	err := r.Parse(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	r.Run()
}
