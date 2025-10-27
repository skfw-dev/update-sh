//go:build linux

package ubuntu

import (
	"update-sh/nextgen/internal/distros"
	"update-sh/nextgen/internal/distros/common"
)

type FlatpakManager struct {
	common.Common
	distros.Context
}

func NewFlatpakManager(ctx distros.Context) *FlatpakManager {
	return &FlatpakManager{
		Context: ctx,
	}
}

func (s *FlatpakManager) IsSupported() bool {
	return false
}

func (s *FlatpakManager) Update() error {
	return distros.ErrNotSupported
}

func (s *FlatpakManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *FlatpakManager) Clean() error {
	return distros.ErrNotSupported
}
