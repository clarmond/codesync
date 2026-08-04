package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	errColor    = color.New(color.FgHiRed)
	greenColor  = color.New(color.FgHiGreen)
	yellowColor = color.New(color.FgHiYellow)
	whiteColor  = color.New(color.FgHiWhite)
	cyanColor   = color.New(color.FgHiCyan)
)

func logSuccess(format string, a ...any) {
	greenColor.Println(fmt.Sprintf(format, a...))
}

func logFail(format string, a ...any) {
	errColor.Println(fmt.Sprintf(format, a...))
}

func logInfo(format string, a ...any) {
	yellowColor.Println(fmt.Sprintf(format, a...))
}

func logIgnore(format string, a ...any) {
	cyanColor.Println(fmt.Sprintf(format, a...))
}

func fatal(format string, a ...any) {
	errColor.Println(fmt.Sprintf(format, a...))
	os.Exit(1)
}
