package config

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/dmawardi/Go-Template/ent"
	"github.com/gorilla/sessions"
)

type AppConfig struct {
	// TemplateCache map[string]*template.Template
	// UseCache      bool
	InProduction bool
	Ctx          context.Context
	DbClient     *ent.Client
	Session      *sessions.CookieStore
	RBEnforcer   *casbin.Enforcer
}
