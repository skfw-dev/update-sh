//go:build !windows

package windows

import (
	"update-sh/nextgen/internal/distros"
	"update-sh/nextgen/internal/distros/common"
)

type WinGetManager struct {
	common.Common
	distros.Context
}

func NewWinGetManager(ctx distros.Context) *WinGetManager {
	return &WinGetManager{
		Context: ctx,
	}
}

func (s *WinGetManager) IsSupported() bool {
	return false
}

func (s *WinGetManager) Update() error {
	return distros.ErrNotSupported
}

func (s *WinGetManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *WinGetManager) Clean() error {
	return distros.ErrNotSupported
}
