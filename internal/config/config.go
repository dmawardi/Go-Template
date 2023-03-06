package config

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type AppConfig struct {
	// TemplateCache map[string]*template.Template
	// UseCache      bool
	InProduction bool
	Ctx          context.Context
	DbClient     *gorm.DB
	Session      *sessions.CookieStore
	RBEnforcer   *casbin.Enforcer
}
