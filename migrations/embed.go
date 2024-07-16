package migrations

import (
	"embed"
)

// Migrations contains the migrations files as embeds
//
//go:embed *.sql
var Migrations embed.FS
