package esi

import (
	"strconv"
	"strings"

	"github.com/kava-forge/eve-alts/lib/errors"
)

type CharacterData struct {
	Name   string `mapstructure:"name"`
	FullID string `mapstructure:"sub"`
	RealID int64  `mapstructure:"-"`
}

func (d *CharacterData) fillRealID() (err error) {
	tmp, _ := strings.CutPrefix(d.FullID, "CHARACTER:EVE:")
	d.RealID, err = strconv.ParseInt(tmp, 10, 64)
	return errors.Wrap(err, "could not convert character id to int64")
}

type CharacterPublicData struct {
	AllianceID    int64  `json:"alliance_id"`
	CorporationID int64  `json:"corporation_id"`
	Name          string `json:"name"`
}

type CharacterPortait struct {
	XLarge string `json:"px512x512"`
	Large  string `json:"px256x256"`
	Medium string `json:"px128x128"`
	Small  string `json:"px64x64"`
}
