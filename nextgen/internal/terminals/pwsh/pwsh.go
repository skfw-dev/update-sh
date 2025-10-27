//go:build !linux && !windows

package pwsh

import (
	"update-sh/nextgen/internal/terminals"
	"update-sh/nextgen/internal/terminals/common"
)

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
