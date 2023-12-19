package main

import (
	"os"

	"github.com/taylormonacelli/boldwanderer"
)

func main() {
	code := boldwanderer.Execute()
	os.Exit(code)
}
