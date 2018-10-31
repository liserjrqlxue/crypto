package AES

import (
	"crypto/aes"
	"crypto/cipher"
)

// []byte length 16/24/32 ~ AES 128/192/256
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func Encode(inputData []byte, codeKey []byte) ([]byte,error){
	c,e:=aes.NewCipher(codeKey)
	if e!=nil{
		return nil,e
	}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	outputData := make([]byte, len(inputData))
	cfb.XORKeyStream(outputData, inputData)
	return outputData,e
}

func Decode(inputData []byte, codeKey []byte) ([]byte,error){
	c,e:=aes.NewCipher(codeKey)
	if e!=nil{
		return nil,e
	}
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	outputData := make([]byte, len(inputData))
	cfbdec.XORKeyStream(outputData, inputData)
	return outputData,e
}
