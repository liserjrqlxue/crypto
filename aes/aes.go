package AES

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"

	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

// []byte length 16/24/32 ~ AES 128/192/256
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func Encode(inputData []byte, codeKey []byte) ([]byte, error) {
	c, e := aes.NewCipher(codeKey)
	if e != nil {
		return nil, e
	}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	outputData := make([]byte, len(inputData))
	cfb.XORKeyStream(outputData, inputData)
	return outputData, e
}

func Decode(inputData []byte, codeKey []byte) ([]byte, error) {
	var c, e = aes.NewCipher(codeKey)
	if e != nil {
		return nil, e
	}
	var cfbdec = cipher.NewCFBDecrypter(c, commonIV)
	var outputData = make([]byte, len(inputData))
	cfbdec.XORKeyStream(outputData, inputData)
	return outputData, e
}

func DecodeFile(fileName string, codeKey []byte) []byte {
	return simpleUtil.HandleError(
		Decode(
			simpleUtil.HandleError(ioutil.ReadFile(fileName)).([]byte),
			codeKey,
		),
	).([]byte)
}

func EncodeFile(fileName string, codeKey []byte) []byte {
	return simpleUtil.HandleError(
		Encode(
			simpleUtil.HandleError(ioutil.ReadFile(fileName)).([]byte),
			codeKey,
		),
	).([]byte)
}

func Decode2File(fileName string, data, codeKey []byte) {
	var file = osUtil.Create(fileName)
	defer simpleUtil.DeferClose(file)
	var c = simpleUtil.HandleError(
		file.Write(
			simpleUtil.HandleError(Decode(data, codeKey)).([]byte),
		),
	)
	fmt.Printf("write %d byte to %s\n", c, fileName)
}

func Encode2File(fileName string, data, codeKey []byte) {
	var file = osUtil.Create(fileName)
	defer simpleUtil.DeferClose(file)
	var c = simpleUtil.HandleError(
		file.Write(
			simpleUtil.HandleError(Encode(data, codeKey)).([]byte),
		),
	)
	fmt.Printf("write %d byte to %s\n", c, fileName)
}

func DecodeFile2File(input, output string, codeKey []byte) {
	var f = osUtil.Create(output)
	defer simpleUtil.DeferClose(f)
	var c = simpleUtil.HandleError(
		f.Write(
			simpleUtil.HandleError(
				Decode(
					simpleUtil.HandleError(ioutil.ReadFile(input)).([]byte),
					codeKey,
				)).([]byte),
		),
	)
	fmt.Printf("write %d byte to %s\n", c, output)
}

func EncodeFile2File(input, output string, codeKey []byte) {
	var f = osUtil.Create(output)
	defer simpleUtil.DeferClose(f)
	var c = simpleUtil.HandleError(
		f.Write(
			simpleUtil.HandleError(
				Encode(
					simpleUtil.HandleError(ioutil.ReadFile(input)).([]byte),
					codeKey,
				)).([]byte),
		),
	)
	fmt.Printf("write %d byte to %s\n", c, output)
}
