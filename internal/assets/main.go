package assets

import "github.com/gobuffalo/packr/v2"

//go:generate packr2
var Migrations = packr.New("migrations", "./migrations")
