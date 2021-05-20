package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	csvCtrl "github.com/danieeelfr/csv_extractor_cli/internal/controller/csv_parser"
	empCtrl "github.com/danieeelfr/csv_extractor_cli/internal/controller/employee"
	fileCtrl "github.com/danieeelfr/csv_extractor_cli/internal/controller/file_parser"
	printCtrl "github.com/danieeelfr/csv_extractor_cli/internal/controller/printer"
	"github.com/danieeelfr/csv_extractor_cli/internal/domain"
	"github.com/dixonwille/wmenu/v5"
	"github.com/pkg/profile"
)

const (
	welcomeMsg = "Welcome to the CSV_EXTRACTOR_CLI system!"
)

var (
	printer               printCtrl.CLIPrinter
	employee              empCtrl.Employee
	csvParser             csvCtrl.CSVParser
	fileParser            fileCtrl.FileParser
	summary               domain.Summary
	inputPath, outputPath string
)

func main() {
	prof := profile.Start(profile.MemProfile)
	checkFlags()
	startControllers()
	printer.IntroScreen()

	menu := startMenu()

	menu.Action(func(opts []wmenu.Opt) error {
		if opts[0].Value == "yes" {
			defer buildSummary()
			defer timeTrack(time.Now(), "extraction proccess")
			defer prof.Stop()
			start()
		} else {
			printer.PrintInfo("proccess not started")
		}
		return nil
	})

	err := menu.Run()
	if err != nil {
		printer.PrintError(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func start() {
	fileNames, err := fileParser.SearchFiles(inputPath)
	if err != nil {
		printer.PrintError(err.Error())
		os.Exit(1)
	}

	printer.PrintInfo(fmt.Sprintf("%v files found inside the %s folder!", len(fileNames), inputPath))

	fmt.Println()
	allEmployees := make([]*domain.Employee, 0)
	for i, fn := range fileNames {
		allEmployees = append(allEmployees, process(i, fn)...)
	}

	printer.PrintInfo("Executing the dedupe process...")
	deduplicated := employee.Dedupe(allEmployees)

	fmt.Println()
	printer.PrintInfo("Spliting valid and invalid records...")
	v, i := employee.SplitValidAndInvalidRecords(deduplicated)
	printer.PrintDebug(fmt.Sprintf("found valid: %v invalid: %v records", len(v), len(i)))

	fmt.Println()
	printer.PrintInfo("Creating file with valid records...")
	err = fileParser.Create(employee.ToCSVString(v), outputPath, "valid.csv")
	if err != nil {
		printer.PrintError(err.Error())
		os.Exit(1)
	}

	printer.PrintInfo("Creating file with invalid records...")
	err = fileParser.Create(employee.ToCSVString(i), outputPath, "invalid.csv")
	if err != nil {
		printer.PrintError(err.Error())
		os.Exit(1)
	}

	fmt.Println()
	printer.PrintInfo("Finished, the output files was generated with success!")
	fmt.Println()

	fetchSummary(len(fileNames), len(allEmployees), len(deduplicated), len(v), len(i))
}

func fetchSummary(files, records, deduplicated, v, i int) {
	summary.Files = files
	summary.Records = records
	summary.Duplicated = records - deduplicated
	summary.Valid = v
	summary.Invalid = i
}

func buildSummary() {
	s := make([][]string, 0)
	h := []string{"Files", "Records", "Duplicated", "Valid", "Invalid", "Execution time"}
	r := summary.ToString()
	s = append(s, h, r)
	printer.PrintSummary(s)
}

func startMenu() *wmenu.Menu {
	menu := wmenu.NewMenu(fmt.Sprintf("\ninput folder: %s\noutput path: %s \n\nLet's start?", inputPath, outputPath))
	menu.IsYesNo(1)
	menu.LoopOnInvalid()
	return menu
}

func process(i int, fn string) []*domain.Employee {
	printer.PrintInfo(fmt.Sprintf("File (%v) %s - processing...", i+1, fn))

	data, err := csvParser.ToCSV(inputPath, outputPath, fn)
	if err != nil {
		printer.PrintError(err.Error())
	}

	fileEmployees := employee.FetchFromStringCSV(data)

	return fileEmployees
}

func startControllers() {
	printer = printCtrl.NewCLIPrinter()
	employee = empCtrl.NewEmployee()
	csvParser = csvCtrl.NewCSVParser()
	fileParser = fileCtrl.NewFileParser()
}

func checkFlags() {
	flag.StringVar(&inputPath, "input_path", "./testdata", "")
	flag.StringVar(&outputPath, "output_path", "./output", "")
	flag.Parse()
}

func timeTrack(start time.Time, name string) {
	summary.ExecutionTime = fmt.Sprint(time.Since(start))
}
