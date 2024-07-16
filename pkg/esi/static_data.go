package esi

import (
	_ "embed"
	"io"

	"github.com/kava-forge/eve-alts/lib/errors"
)

//go:embed staticdata.sqlite
var staticDataDB []byte

func WriteDatabase(to io.Writer) error {
	_, err := to.Write(staticDataDB)
	return errors.Wrap(err, "could not write static data to disk")
}
