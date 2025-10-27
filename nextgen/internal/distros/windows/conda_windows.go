//go:build windows

package windows

import (
	"update-sh/nextgen/internal/distros"
	"update-sh/nextgen/internal/distros/common"
)

type CondaManager struct {
	common.Common
	distros.Context
}

func NewCondaManager(ctx distros.Context) *CondaManager {
	return &CondaManager{
		Context: ctx,
	}
}

func (s *CondaManager) IsSupported() bool {
	return false
}

func (s *CondaManager) Update() error {
	return distros.ErrNotSupported
}

func (s *CondaManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *CondaManager) Clean() error {
	return distros.ErrNotSupported
}
