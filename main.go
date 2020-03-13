package main

import (
	"encoding/csv"
	"flag"
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

	f1 := flag.String("f1", "ARM.csv", "first file to be merged with the second file")
	f2 := flag.String("f2", "X86.csv", "second file to be merged with the first file")

	libraryMapArm := getMap(*f1)
	libraryMapX86 := getMap(*f2)

	file, _ := os.Create("fn2fn.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	_ = writer.Write([]string{"LibraryName", "BinaryName", "FunctionName",
		"Version", "Architecture", "Compiler", "Optimization", "Obfuscation", "EdgeCoverage",
		"Version", "Architecture", "Compiler", "Optimization", "Obfuscation", "EdgeCoverage",
	})
	defer writer.Flush()
	for libraryName, binaryMapArm := range libraryMapArm {
		// find all binaries in the same library
		binaryMapX86, ok := libraryMapX86[libraryName]
		if !ok {
			continue
		}
		for binaryName, functionMapArm := range binaryMapArm {
			// find all function in the same binary
			functionMapX86, ok := binaryMapX86[binaryName]
			if !ok {
				continue
			}
			for fnName, fnsArm := range functionMapArm {
				fnsX86, ok := functionMapX86[fnName]
				if !ok {
					continue
				}
				for i := range fnsArm {
					for j := range fnsX86 {
						_ = writer.Write([]string{fnsArm[i].LibraryName, fnsArm[i].BinaryName, fnsArm[i].FunctionName,
							fnsArm[i].Version, fnsArm[i].Architecture, fnsArm[i].Compiler, fnsArm[i].Optimization, fnsArm[i].Obfuscation, fnsArm[i].EdgeCoverage,
							fnsX86[j].Version, fnsX86[j].Architecture, fnsX86[j].Compiler, fnsX86[j].Optimization, fnsX86[j].Obfuscation, fnsX86[j].EdgeCoverage,
						})
					}
				}
				writer.Flush()
			}
		}
	}
}

func getMap(fileName string) map[string]map[string]map[string][]*csvRow {
	functionMap := make(map[string]map[string]map[string][]*csvRow, 1<<10)

	f, _ := os.Open(fileName)
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
		if functionMap[row.LibraryName][row.BinaryName] == nil {
			functionMap[row.LibraryName][row.BinaryName] = map[string][]*csvRow{}
		}
		if functionMap[row.LibraryName][row.BinaryName][row.FunctionName] == nil {
			functionMap[row.LibraryName][row.BinaryName][row.FunctionName] = []*csvRow{}
		}
		functionMap[row.LibraryName][row.BinaryName][row.FunctionName] = append(functionMap[row.LibraryName][row.BinaryName][row.FunctionName], row)
	}
	return functionMap
}
