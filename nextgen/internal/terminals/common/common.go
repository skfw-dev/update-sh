package common

type Common interface {
	IsSupported() bool
	Upgrade() error
}
