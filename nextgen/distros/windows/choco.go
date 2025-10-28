//go:build !windows

package windows

import (
	"update-sh/unstable/nextgen/distros"
	"update-sh/unstable/nextgen/distros/common"
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
