package controller

import (
	"github.com/OchiengEd/doorman-operator/pkg/controller/doorman"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, doorman.Add)
}
