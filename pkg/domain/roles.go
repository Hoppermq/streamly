package domain

type PlatformRole string

const (
	OwnerRole PlatformRole = "owner"
	AdminRole PlatformRole = "admin"
	UserRole  PlatformRole = "user"
)
