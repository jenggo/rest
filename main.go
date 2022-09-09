package main

import (
	"fmt"
	"os"

	"rest/app"
	"rest/vars"
)

func main() {
	if err := app.Run(&vars.Options{}); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
}
