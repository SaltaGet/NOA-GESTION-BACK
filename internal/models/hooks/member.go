package hooks

import (
	"context"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/aarondl/sqlboiler/v4/boil"
)

func InitMemberHooks() {
    tenant.AddMemberHook(boil.BeforeInsertHook, hashMemberPassword)
    tenant.AddMemberHook(boil.BeforeUpdateHook, hashMemberPassword)
	}

func hashMemberPassword(ctx context.Context, exec boil.ContextExecutor, o *tenant.Member) error {
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