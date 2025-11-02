package key

type ctxKey string

const (
    AppKey    ctxKey = "app"
    TenantDBKey ctxKey = "tenant_db"
)
