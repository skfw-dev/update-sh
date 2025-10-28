//go:build !linux

package ubuntu

import (
	"update-sh/unstable/nextgen/distros"
	"update-sh/unstable/nextgen/distros/common"
)

type APTManager struct {
	common.Common
	distros.Context
}

func NewAPTManager(ctx distros.Context) *APTManager {
	return &APTManager{
		Context: ctx,
	}
}

func (s *APTManager) IsSupported() bool {
	return false
}

func (s *APTManager) Update() error {
	return distros.ErrNotSupported
}

func (s *APTManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *APTManager) Clean() error {
	return distros.ErrNotSupported
}
