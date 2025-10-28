//go:build windows

package windows

import (
	"update-sh/unstable/nextgen/distros"
	"update-sh/unstable/nextgen/distros/common"
)

type ScoopManager struct {
	common.Common
	distros.Context
}

func NewScoopManager(ctx distros.Context) *ScoopManager {
	return &ScoopManager{
		Context: ctx,
	}
}

func (s *ScoopManager) IsSupported() bool {
	return false
}

func (s *ScoopManager) Update() error {
	return distros.ErrNotSupported
}

func (s *ScoopManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *ScoopManager) Clean() error {
	return distros.ErrNotSupported
}
