package auth

import (
	"adb-backup/internal/auth/ldap"
)

type Type int

const (
	None Type = iota
	LDAP
)

var Providers = map[Type]Provider{
	LDAP: &ldap.Provider,
}

var typeNames = map[Type]string{
	None: "none",
	LDAP: "ldap",
}

func (t Type) String() string {
	return typeNames[t]
}

func (t Type) Int() int {
	return int(t)
}

type Provider interface {
	IsReady() bool
	Authenticate(username, password string) (bool, error)
}
