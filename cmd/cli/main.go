package main

import (
	"zero-backend/modules/cli/command"
)

func main() {
	ctx := wireCLIContext()

	if err := command.Execute(ctx); err != nil {
		panic(err)
	}
}
