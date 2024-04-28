package config

import ut "github.com/Rhaqim/creds/internal/utils"

// Postgres
var (
	Database   = ut.Env("DATABASE", "postgres")
	PgHost     = ut.Env("DB_HOST", "localhost")
	PgPort     = ut.Env("DB_PORT", "5432")
	PgUser     = ut.Env("DB_USER", "postgres")
	PgPassword = ut.Env("DB_PASSWORD", "postgres")
	SSLMode    = ut.Env("SSL_MODE", "disable")
	PgUrl      = ut.Env("PG_URL", "postgres://"+PgUser+":"+PgPassword+"@"+PgHost+":"+PgPort+"/"+Database+"?sslmode="+SSLMode)
)
