//go:build linux

package zsh

import (
	"update-sh/nextgen/terminals"
	"update-sh/nextgen/terminals/common"
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
