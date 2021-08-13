package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
	"github.com/xuri/excelize/v2"

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
	skipWarn = flag.String(
		"skipWarn",
		"",
		"skip warn of columns index (0-based), comma as sep",
	)
)

func main() {
	flag.Parse()
	if *input == "" {
		flag.Usage()
		fmt.Println("-input is required!")
		os.Exit(1)
	}
	if *output == "" {
		if *prefix == "" {
			*prefix = *input
		}
		if *suffix == "" {
			switch *outputType {
			case "json":
				*suffix = ".json.aes"
			case "txt":
				*suffix = ".txt.aes"
			}
		}
	}

	var skip = make(map[int]bool)
	if *skipWarn != "" {
		for _, index := range strings.Split(*skipWarn, ",") {
			var i, err = strconv.Atoi(index)
			simpleUtil.CheckErr(err, "can not parse "+*skipWarn)
			skip[i] = true
		}
	}
	var codeKeyBytes = []byte(*codeKey)

	switch *inputType {
	case "xlsx":
		var inputFh, err = excelize.OpenFile(*input)
		simpleUtil.CheckErr(err)
		//var inputFh = simpleUtil.HandleError(excelize.OpenFile(*input)).(*excelize.File)
		//fmt.Printf("%+v\n",inputFh.GetSheetMap())
		for _, sheet := range inputFh.GetSheetMap() {
			if *sheetName != "" && *sheetName != sheet {
				fmt.Printf("skip sheet:[%s]\n", sheet)
				continue
			}
			fmt.Printf("encode sheet:[%s]\n", *sheetName)
			var rows = simpleUtil.HandleError(inputFh.GetRows(sheet)).([][]string)
			fmt.Printf("rows:\t%d\n", len(rows))
			var outputFile = *prefix + *output + "." + sheet + *suffix
			switch *outputType {
			case "json":
				var d []byte
				if *key == "" {
					var data, _ = simpleUtil.Slice2MapArray(rows)
					d = simpleUtil.HandleError(json.MarshalIndent(data, "", "  ")).([]byte)
					fmt.Println("keys:\t%d\n", len(data))
				} else {
					var data, _ = simpleUtil.Slice2MapMapArrayMerge1(rows, *key, *mergeSep, skip)
					d = simpleUtil.HandleError(json.MarshalIndent(data, "", "  ")).([]byte)
					fmt.Printf("keys:\t%d\n", len(data))
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
				d = simpleUtil.HandleError(json.MarshalIndent(data, "", "  ")).([]byte)
			} else {
				var data, _ = simpleUtil.Slice2MapMapArrayMerge1(textUtil.File2Slice(*input, *txtSep), *key, *mergeSep, skip)
				d = simpleUtil.HandleError(json.MarshalIndent(data, "", "  ")).([]byte)
			}
			AES.Encode2File(outputFile, d, codeKeyBytes)
		case "txt":
			AES.EncodeFile2File(*input, outputFile, codeKeyBytes)
		}
	}
}
