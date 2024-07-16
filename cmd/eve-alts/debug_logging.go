//go:build debug

package main

import (
	"fmt"
	stdLog "log"
	"os"
	"path"
)

func maybeEnableFyneLog(home string) {
	fynelog := path.Join(home, "fyne.log")
	w, err := os.Create(fynelog)
	if err != nil {
		panic(err)
	}
	stdLog.SetOutput(w)
	_ = stdLog.Output(1, fmt.Sprintf("os.Args=%v", os.Args))
}
