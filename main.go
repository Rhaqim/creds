package main

import (
	"github.com/Rhaqim/creds/internal/web/router"
)

func main() {
	err := router.Init()
	if err != nil {
		panic(err)
	}
}
