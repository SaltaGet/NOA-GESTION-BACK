package hooks

import (
	"context"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/aarondl/sqlboiler/v4/boil"
)

func InitAdminHooks() {
    // Registramos la función para antes de insertar y antes de actualizar
    master.AddAdminHook(boil.BeforeInsertHook, hashAdminPassword)
    master.AddAdminHook(boil.BeforeUpdateHook, hashAdminPassword)
	}

func hashAdminPassword(ctx context.Context, exec boil.ContextExecutor, o *master.Admin) error {
    if o.Password == "" {
        return nil
    }

    if strings.HasPrefix(o.Password, "$argon2id$") {
        return nil
    }
    
    hash, err := utils.HashPassword(o.Password)
    if err != nil {
        return err
    }
    o.Password = string(hash)
    
    return nil
}