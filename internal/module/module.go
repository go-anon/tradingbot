package module

import (
	"context"

	"github.com/go-anon/simple/logs"
)

type Inter interface {
	// Name returns the module name
	Name() string
	// Start starts the module
	Start(ctx context.Context)
	// Stop stops the module
	Stop()
	// Panic recovers from a panic
	Panic()
	// Clean cleans the module
	Clean()
}

type Module struct {
	name string
}

func New(name string) *Module {
	return &Module{
		name: name,
	}
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Stop() {
	logs.Warn("Received shutdown signal, ( %s ) performing graceful shutdown...", m.Name())
}

func (m *Module) Panic() {
	logs.Warn("Received panic, ( %s ) performing graceful recover...", m.Name())
}

func (m *Module) Clean() {
	logs.Info("Module ( %s ) is clean...", m.Name())
}
