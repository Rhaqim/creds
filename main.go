package main

import (
	"github.com/Rhaqim/creds/internal/database"
	"github.com/Rhaqim/creds/internal/models/migration"
	"github.com/Rhaqim/creds/internal/web/router"
)

func main() {
	database.Init()

	migration.Migrate()

	err := router.Init()
	if err != nil {
		panic(err)
	}

}
