//go:build linux

package ubuntu

import (
	"update-sh/nextgen/internal/distros"
	"update-sh/nextgen/internal/distros/common"
)

type AptManager struct {
	common.Common
	distros.Context
}

func NewAptManager(ctx distros.Context) *AptManager {
	return &AptManager{
		Context: ctx,
	}
}

func (s *AptManager) IsSupported() bool {
	return false
}

func (s *AptManager) Update() error {
	return distros.ErrNotSupported
}

func (s *AptManager) Upgrade() error {
	return distros.ErrNotSupported
}

func (s *AptManager) Clean() error {
	return distros.ErrNotSupported
}
