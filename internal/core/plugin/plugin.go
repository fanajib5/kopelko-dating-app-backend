package plugin

import (
	"context"
	"fmt"
	"log/slog"

	"kopelko-dating-app-backend/internal/core/hook"
	"kopelko-dating-app-backend/internal/platform/config"
	"kopelko-dating-app-backend/internal/platform/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// AppContext provides platform dependencies to plugins.
type AppContext struct {
	Context    context.Context
	DBPool     *pgxpool.Pool
	Transactor database.Transactor
	Hooks      hook.HookManager
	Config     *config.Config
	Logger     *slog.Logger
}

// Plugin defines the standard lifecycle contract for any plugin in Kopelko.
type Plugin interface {
	Name() string
	Version() string
	Init(appCtx *AppContext) error
	RegisterRoutes(apiGroup *echo.Group, authMiddleware echo.MiddlewareFunc)
}

// PluginManager coordinates registered plugins and executes their lifecycle.
type PluginManager struct {
	plugins []Plugin
	appCtx  *AppContext
}

// NewPluginManager instantiates a new manager with AppContext.
func NewPluginManager(appCtx *AppContext) *PluginManager {
	return &PluginManager{
		plugins: make([]Plugin, 0),
		appCtx:  appCtx,
	}
}

// Register adds a plugin to the manager.
func (m *PluginManager) Register(p Plugin) {
	m.plugins = append(m.plugins, p)
}

// InitAll initializes all registered plugins in sequence.
func (m *PluginManager) InitAll() error {
	for _, p := range m.plugins {
		m.appCtx.Logger.Info("Initializing plugin...", "plugin", p.Name(), "version", p.Version())
		if err := p.Init(m.appCtx); err != nil {
			return fmt.Errorf("failed to initialize plugin %s: %w", p.Name(), err)
		}
	}
	return nil
}

// RegisterAllRoutes registers HTTP endpoints for all registered plugins.
func (m *PluginManager) RegisterAllRoutes(apiGroup *echo.Group, authMiddleware echo.MiddlewareFunc) {
	for _, p := range m.plugins {
		p.RegisterRoutes(apiGroup, authMiddleware)
	}
}
