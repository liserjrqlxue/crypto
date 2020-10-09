package main

import (
	"flag"

	AES "github.com/liserjrqlxue/crypto/aes"
)

var (
	input = flag.String(
		"input",
		"",
		"input data",
	)
	output = flag.String(
		"output",
		"",
		"output file name, default is -input.decode",
	)
	codeKey = flag.String(
		"codeKey",
		"c3d112d6a47a0a04aad2b9d2d2cad266",
		"codeKey for aes",
	)
)

func main() {
	flag.Parse()
	if *output == "" {
		*output = *input + ".decode"
	}
	var codeKeyBytes = []byte(*codeKey)

	AES.DecodeFile2File(*input, *output, codeKeyBytes)
}
