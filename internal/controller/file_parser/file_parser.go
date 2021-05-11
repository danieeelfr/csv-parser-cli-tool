package file_parser

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type FileParser interface {
	SearchFiles(path string) ([]string, error)
	Create(data [][]string, path, filename string) error
}

type fileParser struct {
}

func NewFileParser() FileParser {
	return new(fileParser)
}

func (ref *fileParser) SearchFiles(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
		// os.Exit(1)
	}

	names := make([]string, 0)
	for _, f := range files {
		if strings.Contains(f.Name(), ".csv") {
			names = append(names, f.Name())
		}
	}

	return names, nil
}

func (ref *fileParser) Create(data [][]string, path, filename string) error {
	f, err := os.Create(fmt.Sprintf(`%s/%s`, path, filename))
	if err != nil {
		fmt.Println(err)
		return err
		//os.Exit(1)
	}

	writer := csv.NewWriter(f)
	writer.UseCRLF = true

	err = writer.WriteAll(data)
	if err != nil {
		fmt.Println(err)
	}

	return nil

}
