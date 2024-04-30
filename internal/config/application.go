package config

import ut "github.com/Rhaqim/creds/internal/utils"

var (
	// Domain is the domain of the application
	Domain      = ut.Env("DOMAIN", "localhost")
	FrontendURL = ut.Env("FRONTEND_URL", "http://localhost:3000")
)
