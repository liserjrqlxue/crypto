package main

import (
	"encoding/json"
	"flag"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"

	AES "github.com/liserjrqlxue/crypto/aes"
)

var (
	input = flag.String(
		"input",
		"",
		"input",
	)
	inputType = flag.String(
		"type",
		"xlsx",
		"input data type",
	)
	output = flag.String(
		"output",
		"",
		"output file name",
	)
	outputType = flag.String(
		"outputType",
		"json",
		"output data type")
	prefix = flag.String(
		"prefix",
		"",
		"output prefix",
	)
	suffix = flag.String(
		"suffix",
		"",
		"output suffix",
	)
	sheetName = flag.String(
		"sheet",
		"",
		"sheet name of xlsx",
	)
	key = flag.String(
		"key",
		"",
		"key of each line/row",
	)
	mergeSep = flag.String(
		"mergeSep",
		"\n",
		"sep of merge",
	)
	txtSep = flag.String(
		"txtSep",
		"\t",
		"sep of txt",
	)
	codeKey = flag.String(
		"codeKey",
		"c3d112d6a47a0a04aad2b9d2d2cad266",
		"codeKey for aes",
	)
)

func main() {
	flag.Parse()
	if *input == "" {
		if *prefix == "" {
			*prefix = *input
		}
		if *suffix == "" {
			switch *outputType {
			case "json":
				*suffix = ".json"
			case "txt":
				*suffix = ".txt"
			}
		}
	}
	var codeKeyBytes = []byte(*codeKey)

	switch *inputType {
	case "xlsx":
		var inputFh = simpleUtil.HandleError(excelize.OpenFile(*input)).(*excelize.File)
		for sheet := range inputFh.Sheet {
			if *sheetName != "" && *sheetName != sheet {
				continue
			}
			var rows = simpleUtil.HandleError(inputFh.GetRows(sheet)).([][]string)
			var outputFile = *prefix + *output + "." + sheet + *suffix
			switch *outputType {
			case "json":
				var d []byte
				if *key == "" {
					var data, _ = simpleUtil.Slice2MapArray(rows)
					d = simpleUtil.HandleError(json.MarshalIndent(data, "", "")).([]byte)
				} else {
					var data, _ = simpleUtil.Slice2MapMapArrayMerge(rows, *key, *mergeSep)
					d = simpleUtil.HandleError(json.MarshalIndent(data, "", "")).([]byte)
				}
				AES.Encode2File(outputFile, d, codeKeyBytes)
			case "txt":
				var lines []string
				for _, row := range rows {
					lines = append(lines, strings.Join(row, *txtSep))
				}
				AES.Encode2File(outputFile, []byte(strings.Join(lines, *mergeSep)), codeKeyBytes)
			}
		}
	case "txt":
		var outputFile = *prefix + *output + *suffix
		switch *outputType {
		case "json":
			var d []byte
			if *key == "" {
				var data, _ = textUtil.File2MapArray(*input, *txtSep, nil)
				d = simpleUtil.HandleError(json.MarshalIndent(data, "", "")).([]byte)
			} else {
				var data, _ = simpleUtil.Slice2MapMapArrayMerge(textUtil.File2Slice(*input, *txtSep), *key, *mergeSep)
				d = simpleUtil.HandleError(json.MarshalIndent(data, "", "")).([]byte)
			}
			AES.Encode2File(outputFile, d, codeKeyBytes)
		case "txt":
			AES.EncodeFile2File(*input, outputFile, codeKeyBytes)
		}
	}
}
