package main

import (
	"decode-string/services"
	"fmt"
)

func main() {
	input1 := "LR"
	input2 := "RL"
	input3 := "LLRR" // 21012
	input4 := "RRLL" // 01210

	fmt.Println("Input:", input1, "Output:", services.DecodeV2(input1))
	fmt.Println("Input:", input2, "Output:", services.DecodeV2(input2))
	fmt.Println("Input:", input3, "Output:", services.DecodeV2(input3))
	fmt.Println("Input:", input4, "Output:", services.DecodeV2(input4))
}
