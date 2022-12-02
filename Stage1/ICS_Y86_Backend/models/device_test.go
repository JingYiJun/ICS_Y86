package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeviceMarshal(t *testing.T) {
	device := Device{
		Stat: 1,
		Mem:  new(Memory),
	}

	device.Mem.Write(60, 123456, 8)
	device.Mem.Write(64, 1234567, 8)
	device.Mem.Write(600, 12345678912345, 8)
	for i := 60; i < 68; i++ {
		fmt.Printf("%v: %v\n", i, device.Mem[i])
	}

	ans, err := json.Marshal(device)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ans))

	var ansPretty bytes.Buffer
	err = json.Indent(&ansPretty, ans, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(ansPretty.String())
}
