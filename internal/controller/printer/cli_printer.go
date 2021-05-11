package printer

import (
	"fmt"

	"github.com/pterm/pterm"
)

const (
	dfcBGHeader = pterm.BgBlue
)

type CLIPrinter interface {
	PrintInfo(text string)
	PrintDebug(text string)
	PrintError(text string)
	PrintSummary(data [][]string)
	IntroScreen()
}

type CLIPrinterConfigs struct {
}

func NewCLIPrinter() CLIPrinter {
	return new(CLIPrinterConfigs)
}

func (ref *CLIPrinterConfigs) PrintInfo(text string) {
	pterm.Info.Println(text)
}

func (ref *CLIPrinterConfigs) PrintDebug(text string) {
	pterm.Debug.Println(text)
}

func (ref *CLIPrinterConfigs) PrintError(text string) {
	pterm.Error.Println(text)
}

func (ref *CLIPrinterConfigs) IntroScreen() {
	fmt.Sprintln()

	logo, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("CSV:", pterm.NewStyle(pterm.FgLightGreen)),
		pterm.NewLettersFromStringWithStyle(":parser", pterm.NewStyle(pterm.FgLightWhite))).
		Srender()

	pterm.DefaultCenter.Print(logo)

	pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithMargin(0).Sprint("WELCOME!"))
}

func (ref *CLIPrinterConfigs) PrintSummary(data [][]string) {
	pterm.Info.Println("Execution summary:")

	td := pterm.TableData{}

	for _, o := range data {
		td = append(td, o)
	}

	pterm.DefaultTable.WithHasHeader().WithData(td).Render()

}

func clear() {
	print("\033[H\033[2J")
}
