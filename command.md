## Get SQL from DB

```bash
pg_dump -U postgres -h localhost -p 5432 -d noa_gestion --schema-only > internal/platform/database/schemas_db/struct_main.sql

pg_dump -U postgres -h localhost -p 5432 -d daniel_daniel --schema-only > internal/platform/database/schemas_db/struct_tenant.sql
```

## generate models with sqlboiler

```bash
sqlboiler psql -c ./sqlboiler.master.toml 

sqlboiler psql -c ./sqlboiler.tenant.toml 
```

docker exec -it postgres-noa psql -U postgres

\l listar db
\c name_db usar db
\dt listar tablas
\q salir