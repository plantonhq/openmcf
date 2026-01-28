package main

import (
	"github.com/plantonhq/openmcf/cmd/openmcf"
	clipanic "github.com/plantonhq/openmcf/internal/cli/panic"
)

func main() {
	finished := new(bool)
	defer clipanic.Handle(finished)
	openmcf.Execute()
	*finished = true
}
