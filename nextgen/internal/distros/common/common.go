package common

type Common interface {
	IsSupported() bool
	Update() error
	Upgrade() error
	Clean() error
}
