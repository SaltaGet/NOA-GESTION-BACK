package migrations

import "embed"

//go:embed all:main
var MainMigrationsFS embed.FS

//go:embed all:tenant
var TenantMigrations embed.FS
