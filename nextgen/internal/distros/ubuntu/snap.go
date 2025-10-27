//go:build !linux

package ubuntu

import (
	"update-sh/nextgen/internal/distros"
	"update-sh/nextgen/internal/distros/common"
)

type SnapManager struct {
	common.Common
	distros.Context
}

func NewSnapManager(ctx distros.Context) *SnapManager {
	return &SnapManager{
		Context: ctx,
	}
}

func (s *SnapManager) IsSupported() bool {
	return false
}

func (s *SnapManager) Update() error {
	return distros.ErrNotSupported
}

func (s *SnapManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *SnapManager) Clean() error {
	return distros.ErrNotSupported
}
