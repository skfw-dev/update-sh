package distro

type DistroImpl interface {
	GetID() string
	GetFamily() string
	GetIDLike() string
	GetPrimaryPackageManager() string
	String() string
}
