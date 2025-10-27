//go:build !windows

package windows

import (
	"update-sh/nextgen/internal/distros"
	"update-sh/nextgen/internal/distros/common"
)

type Manager struct {
	common.Common
	distros.Context
}

func NewManager(ctx distros.Context) *Manager {
	return &Manager{
		Context: ctx,
	}
}

func (s *Manager) IsSupported() bool {
	return false
}

func (s *Manager) Update() error {
	return distros.ErrNotSupported
}

func (s *Manager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *Manager) Clean() error {
	return distros.ErrNotSupported
}
