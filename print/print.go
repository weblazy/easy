package print

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
)

func PrintBanner(name string) {
	myFigure := figure.NewFigure(name, "", true)
	myFigure.Print()
	fmt.Println()
}
