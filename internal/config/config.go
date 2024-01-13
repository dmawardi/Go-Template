package config

import (
	"context"
	"html/template"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type AppConfig struct {
	// TemplateCache map[string]*template.Template
	// UseCache      bool
	InProduction   bool
	Ctx            context.Context
	DbClient       *gorm.DB
	Session        *sessions.CookieStore
	Auth           AuthEnforcer
	AdminTemplates *template.Template
}

type AuthEnforcer struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}
