package display

import "fmt"

const (
	RESET = "\033[0m"
	RED   = "\033[31m"
	GREEN = "\033[32m"
)

func PrintColored(text, color string) {
	fmt.Printf("%s%s%s\n", color, text, RESET)
}

func PrintColoredBlock(texts []string, color string) {
	for _, t := range texts {
		PrintColored(t, color)
	}
}
