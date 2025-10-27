//go:build windows

package zsh

import (
	"update-sh/nextgen/internal/terminals"
	"update-sh/nextgen/internal/terminals/common"
)

// Manager implements the terminal manager for Zsh on Windows.
// While Zsh can be run on Windows via compatibility layers like Cygwin or
// MSYS2, this manager currently marks it as unsupported.
type Manager struct {
	common.Common
	terminals.Context
}

func NewManager(ctx terminals.Context) *Manager {
	return &Manager{
		Context: ctx,
	}
}

func (s *Manager) IsSupported() bool {
	return false
}

func (s *Manager) Upgrade() error {
	return terminals.ErrNotSupported
}
