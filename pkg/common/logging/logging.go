package logging

import (
	"log"

	"github.com/fatih/color"
)

func InitializeSimpleFormat() {
	log.SetFlags(0)
}

func LogSuccess(text string) {
	color.Set(color.FgGreen)
	log.Println(text)
	color.Unset()
}

func LogError(text string) {
	color.Set(color.FgHiRed)
	log.Println(text)
	color.Unset()
}

func LogSuccessf(text string, v ...interface{}) {
	color.Set(color.FgGreen)
	log.Printf(text, v...)
	color.Unset()
}

func LogErrorf(text string, v ...interface{}) {
	color.Set(color.FgHiRed)
	log.Printf(text, v...)
	color.Unset()
}
