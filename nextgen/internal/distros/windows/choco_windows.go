//go:build windows

package windows

import (
	"update-sh/nextgen/internal/distros"
	"update-sh/nextgen/internal/distros/common"
)

type ChocoManager struct {
	common.Common
	distros.Context
}

func NewChocoManager(ctx distros.Context) *ChocoManager {
	return &ChocoManager{
		Context: ctx,
	}
}

func (s *ChocoManager) IsSupported() bool {
	return false
}

func (s *ChocoManager) Update() error {
	return distros.ErrNotSupported
}

func (s *ChocoManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *ChocoManager) Clean() error {
	return distros.ErrNotSupported
}
