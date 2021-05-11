package csv_parser

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/danieeelfr/csv_extractor_cli/internal/controller/file_parser"
	splitCsv "github.com/tolik505/split-csv"
)

type CSVParser interface {
	ToCSV(inputPath, outputPath, fn string) ([][]string, error)
}

type csvParser struct {
}

func NewCSVParser() CSVParser {
	return new(csvParser)
}

func (ref *csvParser) ToCSV(inputPath, outputPath, fn string) ([][]string, error) {

	file, err := os.Stat(fmt.Sprintf(`%s/%s`, inputPath, fn))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := make([][]string, 0)

	if file.Size() > 20000000 {
		splitter := splitCsv.New()
		splitter.WithHeader = true
		splitter.FileChunkSize = 20000000 //in bytes (20MB)
		result, err := splitter.Split(fmt.Sprintf(`%s/%s`, inputPath, fn), inputPath+"/temp")
		defer removeTempChunkFiles(inputPath + "/temp")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, chunk := range result {
			file, err := os.Open(chunk)
			if err != nil {
				return nil, err
			}

			defer file.Close()
			reader := csv.NewReader(file)
			reader.FieldsPerRecord = -1

			rawCSVdata, err := reader.ReadAll()
			out = append(out, rawCSVdata...)

			if err != nil {
				return nil, err
			}

		}
	} else {
		file, err := os.Open(fmt.Sprintf(`%s/%s`, inputPath, fn))
		if err != nil {
			return nil, err
		}

		defer file.Close()

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1

		rawCSVdata, err := reader.ReadAll()
		out = append(out, rawCSVdata...)

		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func removeTempChunkFiles(folder string) {
	fp := file_parser.NewFileParser()

	tempFiles, _ := fp.SearchFiles(folder)

	for _, fn := range tempFiles {
		err := os.Remove(folder + "/" + fn)
		if err != nil {
			fmt.Println(err)
		}
	}
}
