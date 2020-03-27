package binShape

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type binary struct {
	libraryName  string
	version      string
	binaryName   string
	architecture string
	compiler     string
	optimization string
	obfuscation  string
}

type binShapeRow struct {
	functionName string
	strings_     string
	numStrings   string
	constants    string
	numConstants string
	callers      string
	numCallers   string
	callees      string
	libcCallees  string
}

func newBinShapeRow(records []string) *binShapeRow {
	return &binShapeRow{
		functionName: records[0],
		strings_:     records[1],
		numStrings:   records[2],
		constants:    records[3],
		numConstants: records[4],
		callers:      records[5],
		numCallers:   records[6],
		callees:      records[7],
		libcCallees:  records[8],
	}
}

func mergeBinShapeAndVex() {
	//dir := flag.String("dir", "", "dir that contains all BinShape CSV files")
	//f1 := flag.String("f1", "ARM.csv", "VEX ARM CSV file to be merged")
	//f2 := flag.String("f2", "X86.csv", "VEX X86 CSV file to be merged")
	//
	//flag.Parse()
	//
	//binShapeMap := getBinShapeMap(*dir)
	//vexArmMap := getVexMap(*f1)
	//vexX86Map := getVexMap(*f2)
}

func getBinShapeMap(dir string) map[binary][][]string {
	binDirs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Cannot open the root directory `%v`, got error: %v", dir, err)
	}

	functionMap := make(map[binary][][]string, 1<<10)

	for _, fileInfo := range binDirs {
		if !fileInfo.IsDir() {
			rows, err := getAllFunctions(path.Join(dir, fileInfo.Name()))
			if err != nil {
				continue
			}
			// libcurl-7.42.1-libcurl.so.4.3.0-x86-gcc-O0.fileInfo
			ss := strings.Split(fileInfo.Name(), "-")
			binaryName := ss[2]
			binaryName = strings.TrimSuffix(binaryName, ".so")
			binaryName = strings.Split(binaryName, ".so.")[0]
			b := binary{
				libraryName:  ss[0],
				version:      ss[1],
				binaryName:   binaryName,
				architecture: ss[3],
				compiler:     ss[4],
				optimization: ss[5],
				obfuscation:  "",
			}
			functionMap[b] = rows
		}
	}

	return functionMap
}

func getAllFunctions(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Cannot open file %v, got error: %v", fileName, err)
		return nil, nil
	}
	defer f.Close() // this needs to be after the err check

	reader := csv.NewReader(f)
	_, _ = reader.Read() // skip headers
	return reader.ReadAll()
}

type vexRow struct {
	functionName string
	edgeCoverage string
	callWalks    string
}

func getVexMap(fileName string) map[binary][]vexRow {
	functionMap := make(map[binary][]vexRow, 1<<10)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Cannot open file %v, got error: %v", fileName, err)
	}
	defer f.Close() // this needs to be after the err check

	reader := csv.NewReader(f)
	_, _ = reader.Read() // skip headers
	for {
		records, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		binaryName := records[2]
		binaryName = strings.TrimSuffix(binaryName, ".so")
		binaryName = strings.Split(binaryName, ".so.")[0]
		b := binary{
			libraryName:  records[0],
			version:      records[1],
			binaryName:   binaryName,
			architecture: records[3],
			compiler:     records[4],
			optimization: records[5],
			obfuscation:  records[6],
			//functionName: records[7],
			//edgeCoverage: records[8],
			//callWalks:    records[9],
		}

		r := vexRow{
			functionName: records[7],
			edgeCoverage: records[8],
			callWalks:    records[9],
		}

		if functionMap[b] == nil {
			functionMap[b] = []vexRow{}
		}

		functionMap[b] = append(functionMap[b], r)
	}
	return functionMap
}
