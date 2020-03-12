package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type csvRow struct {
	LibraryName  string
	Version      string
	BinaryName   string
	Architecture string
	Compiler     string
	Optimization string
	Obfuscation  string
	FunctionName string
	EdgeCoverage string
	CallWalks    string
}

func newRow(records []string) *csvRow {
	return &csvRow{
		LibraryName:  records[0],
		Version:      records[1],
		BinaryName:   records[2],
		Architecture: records[3],
		Compiler:     records[4],
		Optimization: records[5],
		Obfuscation:  records[6],
		FunctionName: records[7],
		EdgeCoverage: strings.TrimSpace(records[8]),
		CallWalks:    records[9],
	}
}

func main() {
	//library name
	//	--- binary name
	//		--- function name
	//			--- []csvRow

	mapArm := getMap("ARM.csv_1000")
	mapX86 := getMap("X86.csv_1000")
	fmt.Printf("%v, %v", len(mapArm), len(mapX86))
}

func getMap(fileName string) map[string]map[string]map[string][]*csvRow {
	functionMap := make(map[string]map[string]map[string][]*csvRow, 1<<10)

	f, _ := os.Open("ARM_10000.csv")
	defer f.Close() // this needs to be after the err check

	reader := csv.NewReader(f)
	_, _ = reader.Read() // skip headers
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		row := newRow(record)
		if functionMap[row.LibraryName] == nil {
			functionMap[row.LibraryName] = map[string]map[string][]*csvRow{}
		}
		if functionMap[row.LibraryName][row.BinaryName] == nil{
			functionMap[row.LibraryName][row.BinaryName] = map[string][]*csvRow{}
		}
		if functionMap[row.LibraryName][row.BinaryName][row.FunctionName] == nil {
			functionMap[row.LibraryName][row.BinaryName][row.FunctionName] = []*csvRow{}
		}
		functionMap[row.LibraryName][row.BinaryName][row.FunctionName] = append(functionMap[row.LibraryName][row.BinaryName][row.FunctionName], row)
		fmt.Printf("%+v", row.LibraryName)
	}
	return functionMap
}
