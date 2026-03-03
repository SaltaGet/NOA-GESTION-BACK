#!/bin/bash
cd /home/daniel/Documentos/Programacion/Go/Fiber/NOA-GESTION/NOA-GESTION-BACK

for mod_dir in internal/module/*; do
    if [ ! -d "$mod_dir" ] || [ "$(basename "$mod_dir")" = "client" ]; then
        continue
    fi
    
    entity=$(basename "$mod_dir")
    
    # 1. Schemas
    if [ -f "internal/schemas/${entity}.go" ]; then
        echo -e "package repository\n" > "$mod_dir/infrastructure/repository/schemas.go"
        grep -v "^package " "internal/schemas/${entity}.go" >> "$mod_dir/infrastructure/repository/schemas.go"
    fi
    
    # 2. Models
    if [ -f "internal/models/${entity}.go" ]; then
        echo -e "package domain\n" > "$mod_dir/domain/models.go"
        grep -v "^package " "internal/models/${entity}.go" >> "$mod_dir/domain/models.go"
    fi
    
    # 3. Ports 
    if [ -f "internal/ports/${entity}.go" ]; then
        echo -e "package domain\n" > "$mod_dir/domain/ports.go"
        grep -v "^package " "internal/ports/${entity}.go" >> "$mod_dir/domain/ports.go"
    fi
    
    # 4. Services
    if [ -f "internal/services/${entity}.go" ]; then
        echo -e "package application\n" > "$mod_dir/application/services.go"
        grep -v "^package " "internal/services/${entity}.go" >> "$mod_dir/application/services.go"
    fi
    
    # 5. Repositories
    if [ -f "internal/repositories/${entity}.go" ]; then
        echo -e "package repository\n" > "$mod_dir/infrastructure/repository/repository.go"
        grep -v "^package " "internal/repositories/${entity}.go" >> "$mod_dir/infrastructure/repository/repository.go"
    fi
    
    # 6. Controllers
    if [ -f "cmd/api/controllers/${entity}.go" ]; then
        mkdir -p "$mod_dir/infrastructure/handler/http"
        echo -e "package http\n" > "$mod_dir/infrastructure/handler/http/controller.go"
        grep -v "^package " "cmd/api/controllers/${entity}.go" >> "$mod_dir/infrastructure/handler/http/controller.go"
    fi
    
    # 7. Routes
    if [ -f "cmd/api/routes/${entity}.go" ]; then
        mkdir -p "$mod_dir/infrastructure/handler/http"
        echo -e "package http\n" > "$mod_dir/infrastructure/handler/http/routes.go"
        grep -v "^package " "cmd/api/routes/${entity}.go" >> "$mod_dir/infrastructure/handler/http/routes.go"
    fi
    
    echo "Copied files for module: $entity"
done
