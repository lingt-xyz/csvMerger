package fn2fn

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

type vexRow struct {
	libraryName  string
	version      string
	binaryName   string
	architecture string
	compiler     string
	optimization string
	obfuscation  string
	functionName string
	edgeCoverage string
	callWalks    string
}

func newVexRow(records []string) *vexRow {
	binaryName := records[2]
	binaryName = strings.TrimSuffix(binaryName, ".so")
	binaryName = strings.Split(binaryName, ".so.")[0]
	return &vexRow{
		libraryName:  records[0],
		version:      records[1],
		binaryName:   binaryName,
		architecture: records[3],
		compiler:     records[4],
		optimization: records[5],
		obfuscation:  records[6],
		functionName: records[7],
		edgeCoverage: records[8],
		callWalks:    records[9],
	}
}

func MapFunctionsX86AndArm() {
	//library name
	//	--- binary name
	//		--- function name
	//			--- []vexRow

	f1 := flag.String("f1", "ARM.csv", "first file to be merged with the second file")
	f2 := flag.String("f2", "X86.csv", "second file to be merged with the first file")
	flag.Parse()

	libraryMapArm := getVexMap(*f1)
	libraryMapX86 := getVexMap(*f2)

	file, _ := os.Create("fn2fn.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	_ = writer.Write([]string{"libraryName", "binaryName", "functionName",
		"version", "architecture", "compiler", "optimization", "obfuscation", "edgeCoverage",
		"version", "architecture", "compiler", "optimization", "obfuscation", "edgeCoverage",
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
						_ = writer.Write([]string{fnsArm[i].libraryName, fnsArm[i].binaryName, fnsArm[i].functionName,
							fnsArm[i].version, fnsArm[i].architecture, fnsArm[i].compiler, fnsArm[i].optimization, fnsArm[i].obfuscation, fnsArm[i].edgeCoverage,
							fnsX86[j].version, fnsX86[j].architecture, fnsX86[j].compiler, fnsX86[j].optimization, fnsX86[j].obfuscation, fnsX86[j].edgeCoverage,
						})
					}
					writer.Flush()
				}
			}
		}
	}
}

func RemoveEmptyFields() {
	input := flag.String("input", "fn2fn.csv", "CSV file whose empty fields to be removed")
	output := flag.String("output", "X86.csv", "CSV file that has filtered empty fields")
	flag.Parse()

	inputCSV, err := os.Open(*input)
	if err != nil {
		log.Fatalf("Cannot open file %v, got error: %v", *input, err)
	}
	defer inputCSV.Close() // this needs to be after the err check

	outputCSV, err := os.Create(*output)
	if err != nil {
		log.Fatalf("Cannot open file %v, got error: %v", *output, err)
	}
	defer outputCSV.Close()

	writer := csv.NewWriter(outputCSV)
	_ = writer.Write([]string{"libraryName", "binaryName", "functionName",
		"version", "architecture", "compiler", "optimization", "obfuscation", "edgeCoverage",
		"version", "architecture", "compiler", "optimization", "obfuscation", "edgeCoverage",
	})
	defer writer.Flush()

	reader := csv.NewReader(inputCSV)
	_, _ = reader.Read() // skip headers
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[9] != "" && record[14] != "" {
			_ = writer.Write(record)
			writer.Flush()
		}
	}
}

func getVexMap(fileName string) map[string]map[string]map[string][]*vexRow {
	functionMap := make(map[string]map[string]map[string][]*vexRow, 1<<10)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Cannot open file %v, got error: %v", fileName, err)
	}
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
		row := newVexRow(record)
		if functionMap[row.libraryName] == nil {
			functionMap[row.libraryName] = map[string]map[string][]*vexRow{}
		}
		if functionMap[row.libraryName][row.binaryName] == nil {
			functionMap[row.libraryName][row.binaryName] = map[string][]*vexRow{}
		}
		if functionMap[row.libraryName][row.binaryName][row.functionName] == nil {
			functionMap[row.libraryName][row.binaryName][row.functionName] = []*vexRow{}
		}
		functionMap[row.libraryName][row.binaryName][row.functionName] = append(functionMap[row.libraryName][row.binaryName][row.functionName], row)
	}
	return functionMap
}